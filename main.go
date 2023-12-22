package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/udonetsm/server/controller"
)

func main() {
	StartingServer()
}

func StartingServer() {
	mux := mux.NewRouter()
	mux.HandleFunc("/update/number", controller.UpdateNumberController)
	mux.HandleFunc("/update/name", controller.UpdateNameController)
	mux.HandleFunc("/update/numlist", controller.UpdateNumListController)
	mux.HandleFunc("/delete", controller.DeleteController)
	mux.HandleFunc("/info", controller.InfoController)
	mux.HandleFunc("/create", controller.CreateController)

	server := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
		Handler:      mux,
		Addr:         ":8080",
	}
	log.Println("Starting at", time.Now())
	log.Println(server.ListenAndServe())
}
