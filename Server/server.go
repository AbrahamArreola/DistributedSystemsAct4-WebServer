package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadHTML(path string) string {
	html, _ := ioutil.ReadFile(path)

	return string(html)
}

func index(response http.ResponseWriter, request *http.Request) {
	response.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		response,
		loadHTML("../Pages/add_score.html"),
	)
}

func main() {
	host := "127.0.0.1:9000"
	http.HandleFunc("/", index)
	fmt.Println("Servidor corriendo en:", host)
	http.ListenAndServe(host, nil)
}
