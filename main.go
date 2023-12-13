package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/udonetsm/client/models"
)

func main() {
	StartingServer()
}

func StartingServer() {
	mux := mux.NewRouter()
	mux.HandleFunc("/update/number", UpdateNumber)
	server := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
		Handler:      mux,
		Addr:         ":8080",
	}
	log.Println(server.ListenAndServe())
}

func UpdateNumber(w http.ResponseWriter, r *http.Request) {
	a := &models.RequestJSON{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, a)
	if err != nil {
		log.Fatal(err)
	}
	//use a to update contact in db.
}
