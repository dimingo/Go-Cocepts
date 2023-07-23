package main

import (
	"example/sessions"
	"net/http"
	"text/template"
)

var store = sessions.NewCookieStore([]byte("VSLRz84Yfr8EOMHFmK37GYcbm/RHkA9l"))

func Count(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ct := session.Values["countnum"]

	if ct == nil {
		session.Values["countnum"] = 1
	} else {
		session.Values["countnum"] = ct.(int) + 1
	}

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Render the template with the count value
	t, err := template.ParseFiles("count.gtpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, session.Values["countnum"])
}

func main() {
	http.HandleFunc("/count", Count)
	http.ListenAndServe(":9090", nil)
}
