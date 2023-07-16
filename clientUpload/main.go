package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, taregetUrl string) error {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// very important step
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)

	if err != nil {
		fmt.Println("Error writing to buffer")
		return err

	}

	// open file handle
	fh, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error openning the file")

		return err
	}

	defer fh.Close()

	// iocopy
	_, err = io.Copy(fileWriter, fh)

	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//     resp, err := http.Post(targetUrl, contentType, bodyBuf)

	resp, err := http.Post(taregetUrl, contentType, bodyBuf)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil

}

func main() {

	target_url := "http://localhost:9090/upload"
	filename := "./person.txt"
	postFile(filename, target_url)
}
