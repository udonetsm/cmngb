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
	flags.Flags()
}

func BuildServer() {
	mux := mux.NewRouter()
	mux.HandleFunc("/info", controllers.MW(controllers.Info))
	mux.HandleFunc("/get/list", controllers.MW(controllers.Search))
	mux.HandleFunc("/delete", controllers.MW(controllers.Delete))
	mux.HandleFunc("/create", controllers.MW(controllers.Create))
	mux.HandleFunc("/update", controllers.MW(controllers.Update))
	mux.HandleFunc("/get/token", controllers.Auth)
	mux.HandleFunc("/create/user", controllers.NewUser)
	server := &http.Server{
		Handler: mux,
		Addr:    flags.SRVHOST + flags.SRVPORT,
	}

	log.Println("Starting server")
	server.ListenAndServe()
}
