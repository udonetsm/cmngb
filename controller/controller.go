package controller

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
	"github.com/udonetsm/server/use"
)

// This function write in ResponseWriter al of errors, header with http.Status and logging it
func UnificatedErrorCapturing(j *models.Entries, w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
	log.Println("ERROR for", j.Object, "with status", status, err.Error())
}

// Unpack json from request body
// Update column number in database and
// field num of json in database
func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	// Must valid new user Number.
	// User number must contains only numerics
	err := use.MatchNumber(j, c)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	models.UnpackingContact(c, []byte(j.Object))
	err = database.UpdateNumber(j, c.Number)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
	log.Println("Number list updated general number by", j.Number, "new object", j.Object)
}

// this function unpacks json from request
// unpack JSON from this unpacked object
// and call target function from database package.
func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	err := use.MatchName(j, c)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	models.UnpackingContact(c, []byte(j.Object))
	// call database package function using name of upgradable json field in database
	err = database.Update(j, "name", c.Name)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
	log.Println("Number list updated name by", j.Number, "new object", j.Object)

}

// this function updates general number of contct in db and
// same updates field num in json object in database
func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(j.Object))
	err := database.Update(j, "nlist", c.NumberList)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(j.Object))
	log.Println("Number list updated number list by", j.Number, "new object", j.Object)
}

// returns json string of contact from database(field object).
func InfoController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := database.GetInfo(j)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(j.Object))
	log.Println("Get info by", j.Number)
}

// Unpacks json from request body and delete target contact
// Write header bad request and error from database.Delete()
// if error isn't nil
func DeleteController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := database.Delete(j)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j.Object))
	log.Println("Deleted", j.Object)
}

// Unpacks request body, call Create from database
// if Create() returns error, response includes this it
// and status code bad request
func CreateController(w http.ResponseWriter, r *http.Request) {
	j := request(w, r)
	err := matchJsonFieldAndTarget(j)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	err = use.Matching(j)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	err = database.Create(j)
	if err != nil {
		UnificatedErrorCapturing(j, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(j.Object))
	log.Println("Created", j.Object)
}

// Function for matching internal target number
// and JSON object number. If aren't equal write error in ResponseWriter
func matchJsonFieldAndTarget(e *models.Entries) (err error) {
	c := &models.Contact{}
	models.UnpackingContact(c, []byte(e.Object))
	if c.Number != e.Number {
		err = errors.New("TARGET NUMBER AND JSON OBJECT NUMBER AREN'T EQUAL")
	}
	return
}

// This is a local function for unpacking jsons from requests
func request(w http.ResponseWriter, r *http.Request) *models.Entries {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something wrong. Try later"))
	}
	e := &models.Entries{}
	models.Unpacking(e, req)
	return e
}
