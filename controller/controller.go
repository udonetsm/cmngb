package controller

import (
	"io"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
)

func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	err := database.UpdateNumber(j, c.Number)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("[UPDATENUMBER OK]"))
}

func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	err := database.UpdateName(j, c.Name)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("[UPDATENAME OK]"))
}

func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	// j := Request(w, r)
}

func InfoController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	obj := database.GetInfo(j)
	w.Write([]byte("[INFO OK]" + obj))
}

func DeleteController(w http.ResponseWriter, r *http.Request) {
	j := Request(w, r)
	o := database.GetInfo(j)
	if len(o) == 0 {
		w.Write([]byte("[NOT AFFECTED]"))
		return
	}
	err := database.Delete(j)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("[DELETE OK] " + o))
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
	models.Unpacking(o, req)
	return o
}
