package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// session struct
type Session struct {
	ID      string
	IsNew   bool
	values  map[interface{}]interface{}
	name    string
	store   store
	Options *Options
}

// central manager any request passes through here
type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    provider
	maxlifetime int64
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""

	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) sessionStart(w http.ResponseWriter, r *http.Request) (Session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {

		sid := manager.sessionId()
		Session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		Session, _ = manager.provider.SessionRead(sid)

	}

	return
}
func newManager(providerName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[providerName]
	if !ok {
		return nil, fmt.Errorf("session: uknown provide %q (forgottern import?)", providerName)

	}

	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func init() {
	globalSessions = newManager("memory", "gosessionid", 3600)
}

func main() {
	var globalSessions *session.Manager

}
