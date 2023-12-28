package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/use"
)

func Create(w http.ResponseWriter, r *http.Request) {
	log.Println("CREATE")
}

func UpdateGeneralNumber(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	fmt.Println(e)
	// use.Match(e, use.EQAL)
}

func UpdateName(w http.ResponseWriter, r *http.Request) {
}

func UpdateListNumber(w http.ResponseWriter, r *http.Request) {
}

func Info(w http.ResponseWriter, r *http.Request) {
}

func Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE")
}

func MW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := request(w, r)
		use.Match(e, use.ENUM)
		if e.Error != nil {
			fmt.Fprintln(w, e)
			log.Println(e.Error, "for object", e.Object)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// This is a local function for unpacking jsons from requests
func request(w http.ResponseWriter, r *http.Request) *models.Entries {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something wrong. Try later"))
	}
	e := &models.Entries{}
	models.UnpackingEntries(e, req)
	return e
}
