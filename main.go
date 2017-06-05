package main

import (
	"log"
	"net/http"

	lib "github.com/umens/go-url-shortener/lib"
)

func main() {

	http.HandleFunc("/url/", lib.RedirectHandler)
	http.HandleFunc("/r/", lib.RedirectionHandler)
	http.HandleFunc("/shorten", lib.ShorthenHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	log.Printf("About to listen on 8080. Go to https://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", lib.WrapHandlerWithLogging(http.DefaultServeMux)))
}
