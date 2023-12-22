package controller

import (
	"io"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
)

// Unpack json from request body
// Update column number in database and
// field num of json in database
func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	err := database.UpdateNumber(j, c.Number)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
}

// this function unpacks json from request
// unpack JSON from this unpacked object
// and call target function from database package.
func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	// call database package function using name of upgradable json field in database
	err := database.Update(j, "name", c.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
}

// this function updates general number of contct in db and
// same updates field num in json object in database
func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	err := database.Update(j, "nlist", c.NumberList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
}

// returns json string of contact from database(field object).
func InfoController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := database.GetInfo(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(j.Object))
}

// Unpacks json from request body and delete target contact
// Write header bad request and error from database.Delete()
// if error isn't nil
func DeleteController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := database.Delete(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j.Object))
}

// Unpacks request body, call Create from database
// if Create() returns error, response includes this it
// and status code bad request
func CreateController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := database.Create(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(j.Object))
}

// This is a local function for unpacking jsons from requests
func request(w http.ResponseWriter, r *http.Request) *models.Entries {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	o := &models.Entries{}
	models.Unpacking(o, req)
	return o
}
