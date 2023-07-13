package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

//upload logic

func upload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method: ", r.Method)

	if r.Method == "GET" {

		curtime := time.Now().Unix()

		h := md5.New()

		io.WriteString(h, strconv.FormatInt(curtime, 10))

		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("fielUpload.gtpl")
		t.Execute(w, token)

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")

		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		fmt.Println(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return

		}

		defer f.Close()
		io.Copy(f, file)
	}
}
