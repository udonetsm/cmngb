package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
)

func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	req := Request(w, r)
	fmt.Println(req)
}

func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	Request(w, r)

}

func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	Request(w, r)
}

func InfoController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	obj := database.GetInfo(j)
	w.Write([]byte(obj))
}

func DeleteController(w http.ResponseWriter, r *http.Request) {
}

func CreateController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	err := database.Create(j)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("[CREATE OK] " + j.Object))
}

func Request(w http.ResponseWriter, r *http.Request) *models.Entries {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	o := &models.Entries{}
	o.Unpack(req)
	return o
}
