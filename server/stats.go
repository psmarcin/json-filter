package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[STATS] ")
	defer log.SetPrefix("")
	resJSON, err := json.Marshal(rootStatus)
	if err != nil {
		errorResponse(err, w)
		return
	}

	jsonResponse(resJSON, w)
}
