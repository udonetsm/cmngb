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
	mux.HandleFunc("/create", controller.MW(controller.Create))
	mux.HandleFunc("/update/number", controller.MW(controller.UpdateGeneralNumber))
	mux.HandleFunc("/update/name", controller.MW(controller.UpdateName))
	mux.HandleFunc("/update/listnumber", controller.MW(controller.UpdateListNumber))
	mux.HandleFunc("/info", controller.MW(controller.Info))
	mux.HandleFunc("/search", controller.Search)
	mux.HandleFunc("/delete", controller.Delete)

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
