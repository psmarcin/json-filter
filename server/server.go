package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var port = ":8080"

type Status struct {
	Ok        bool      `json:"ok"`
	StartedAt time.Time `json:"startedAt"`
}

type Error struct {
	IsError      bool      `json:"isError"`
	Timestamp    time.Time `json:"timestamp"`
	ErrorMessage string    `json:"error"`
}

var rootStatus = Status{
	Ok:        true,
	StartedAt: time.Now(),
}

const ERR500 = "Internal Error"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resJSON, err := json.Marshal(rootStatus)
	if err != nil {
		error := Error{
			IsError:      true,
			Timestamp:    time.Now(),
			ErrorMessage: string(err.Error()),
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		resJSON, _ := json.Marshal(error)
		fmt.Fprint(w, string(resJSON))
		log.Printf("Error parsing json %v", err)
		return
	}
	log.Printf("Request %v", rootStatus)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(resJSON))
}

func Start() {
	http.HandleFunc("/", handler)

	log.Printf("Starting server at %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
