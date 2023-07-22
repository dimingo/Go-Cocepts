package main

import (
	"example/sessions"
	"fmt"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("VSLRz84Yfr8EOMHFmK37GYcbm/RHkA9l"))

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "session-name")
	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	foo := session.Values["foo"]
	val := session.Values[42]

	// Print the retrieved data to the console
	fmt.Println("foo:", foo)
	fmt.Println("val", val)
}
func main() {
	http.HandleFunc("/", MyHandler)
	err := http.ListenAndServe(":9090", nil) // set listening port
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}

}
