package sessions

import "net/http"

// newCookieFromOptions returns an http.Cookie with the options set.
func newCookieFromOptions(name, value string, options *options) *http.Cookie {

	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
