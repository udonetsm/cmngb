package main

import (
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
	mux.HandleFunc("/search", controllers.MW(controllers.Search))
	mux.HandleFunc("/delete/by/id", controllers.MW(controllers.DeleteById))
	mux.HandleFunc("/create", controllers.Create)
	mux.HandleFunc("/update/number", controllers.MW(controllers.UpdateNumber))
	mux.HandleFunc("/update/name", controllers.MW(controllers.Update))
	mux.HandleFunc("/update/list", controllers.MW(controllers.Update))
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
