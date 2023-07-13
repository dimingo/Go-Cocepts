package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "Hello Mitchel!") // send data to client side
}

func main() {
	http.HandleFunc("/", sayhelloName) // set router

	err := http.ListenAndServe(":9090", nil) // set listening port
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
