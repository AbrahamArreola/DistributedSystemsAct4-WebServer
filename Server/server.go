package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadHTML(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		loadHTML("form.html"),
	)
}

func main() {
	host := "127.0.0.1:9000"
	http.HandleFunc("/", form)
	fmt.Println("Servidor corriendo en:", host)
	http.ListenAndServe(host, nil)
}
