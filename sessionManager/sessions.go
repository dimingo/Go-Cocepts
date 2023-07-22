package sessions

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"
)

// session struct
type Session struct {
	ID      string
	IsNew   bool
	values  map[interface{}]interface{}
	name    string
	store   Store
	Options *Options
}

// Default flashes key.
const flashesKey = "_flash"

// Session --------------------------------------------------------------------

// New session is called by session stores to create new session instance
func NewSession(store Store, name string) *Session {
	return &Session{
		values:  make(map[interface{}]interface{}),
		store:   store,
		name:    name,
		Options: new(Options),
	}
}

// Flashes returns a slice of flash messages from the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) Flashes(vars ...string) []interface{} {
	var flashes []interface{}
	key := flashesKey
	if len(vars) > 0 {
		key = vars[0]
	}
	if v, ok := s.values[key]; ok {
		// Drop the flashes and return it.
		delete(s.values, key)
		flashes = v.([]interface{})
	}
	return flashes
}

// AddFlash adds a flash message to the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) AddFlash(value interface{}, vars ...string) {
	key := flashesKey
	if len(vars) > 0 {
		key = vars[0]
	}
	var flashes []interface{}
	if v, ok := s.values[key]; ok {
		flashes = v.([]interface{})
	}
	s.values[key] = append(flashes, value)
}

// method to save the session
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {

	return s.store.Save(r, w, s)

}

// return the sessopm store used to register the session
func (s *Session) Store() Store {

	return s.store
}

// return the session Name()
func (s *Session) Name() string {
	return s.name
}

// Registry --------------------------------------

// SessionInfo stores a session trackde by registry
type sessionInfo struct {
	s *Session
	e error
}

// contextKey is the type used to store the registry in the context
type contextKey int

// registryKey is the key used to store the registry in th e context
const registryKey contextKey = 0

// GetRegistry returns a registry instance of current request
func GetRegistry(r *http.Request) *Registry {
	var ctx = r.Context()
	registry := ctx.Value(registryKey)
	if registry != nil {
		return registry.(*Registry)
	}
	newRegistry := &Registry{
		request:  r,
		sessions: make(map[string]sessionInfo),
	}
	*r = *r.WithContext(context.WithValue(ctx, registryKey, newRegistry))
	return newRegistry
}

// Registry stores session used during a request
type Registry struct {
	request  http.Request
	sessions map[string]sessionInfo
}

// Get registers and returns a session for a given name and session store
// it returns a new session if there are no sessions registered for the name
func (s *Registry) Get(store Store, name string) (session *Session, err error) {
	if !isCookieNameValid(name) {

		return nil, fmt.Errorf("sessions: invalid character in cookie name: %s", name)
	}

	if info, ok := s.sessions[name]; ok {
		session, err = info.s, info.e
	} else {
		session, err = store.New(s.request, name)
		session.name = name
		s.sessions[name] = sessionInfo{s: session, e: err}
	}
	session.store = store
	return
}

// Save saves all sessions registered on current request
func (s *Registry) Save(w http.ResponseWriter) error {

	var errMulti MultiError
	for name, info := range s.sessions {
		session := info.s
		if session.store == nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: missing store for session %q", name))
		} else if err := session.store.Save(s.request, w, session); err != nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: error saving session %q -- %v", name, err))
		}
	}
	if errMulti != nil {
		return errMulti
	}
	return nil
}

// Heplers ----------------------------------------------------------------------------

func init() {
	gob.Register([]interface{}{})
}

// Save saves all sessions used during the current request.
func Save(r *http.Request, w http.ResponseWriter) error {
	return GetRegistry(r).Save(w)
}

// NewCookie returns an http.Cookie with the options set. It also sets
// the Expires field calculated based on the MaxAge value, for Internet
// Explorer compatibility.
func NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := newCookieFromOptions(name, value, options)
	if options.MaxAge > 0 {
		d := time.Duration(options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if options.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}

// Error ----------------------------------------------------------------------

// MultiError stores multiple errors.
//
// Borrowed from the App Engine SDK.
type MultiError []error

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}

