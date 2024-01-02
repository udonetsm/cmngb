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
	mux.HandleFunc("/info", controllers.Info)
	mux.HandleFunc("/search", controllers.Search)
	mux.HandleFunc("/delete/by/id", controllers.DeleteById)
	mux.HandleFunc("/create", controllers.Create)
	mux.HandleFunc("/update/number", controllers.UpdateNumber)
	mux.HandleFunc("/update/name", controllers.Update)
	mux.HandleFunc("/update/list", controllers.Update)
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
