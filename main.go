package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/udonetsm/cmngb/controllers"
)

func main() {
	BuildServer()
}

func BuildServer() {
	mux := mux.NewRouter()
	mux.HandleFunc("/info", controllers.MW(controllers.Info))
	mux.HandleFunc("/search", controllers.Search)
	mux.HandleFunc("/delete", controllers.MW(controllers.Delete))
	mux.HandleFunc("/create", controllers.MW(controllers.Create))
	mux.HandleFunc("/update", controllers.MW(controllers.Update))
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	log.Println("Starting server")
	server.ListenAndServe()
}
