package server

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request %v", r.URL.Path)
	fmt.Fprintf(w, "Reqest: %s", r.URL.Path)
}

func Start() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
