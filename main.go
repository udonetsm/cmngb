package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/udonetsm/cmngb/controllers"
	"github.com/udonetsm/cmngb/flags"
)

func main() {
	BuildServer()
}

func init() {
	flags.Load()
}

func BuildServer() {
	mux := mux.NewRouter()
	mux.HandleFunc("/inf", controllers.MW(controllers.Info))
	mux.HandleFunc("/lst", controllers.Search)
	mux.HandleFunc("/del", controllers.MW(controllers.Delete))
	mux.HandleFunc("/new", controllers.MW(controllers.Create))
	mux.HandleFunc("/upd", controllers.MW(controllers.Update))
	mux.HandleFunc("/auth", controllers.Auth)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	log.Println("Starting server")
	server.ListenAndServe()
}
