package server

import (
	"fmt"
	"log"
	"net/http"
)

var port = ":8080"

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request %v", r.URL.Path)
	fmt.Fprintf(w, "Reqest: %s", r.URL.Path)
}

func Start() {
	http.HandleFunc("/", handler)
	log.Printf("Starting server at %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
