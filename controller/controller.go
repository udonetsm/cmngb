package controller

import (
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
	"github.com/udonetsm/server/use"
)

// This function write in ResponseWriter al of errors, header with http.Status and logging it
func UnificatedErrorCapturing(e *models.Entries, w http.ResponseWriter, err error, status int) {
	w.Write([]byte(err.Error() + "\n"))
	log.Println("ERROR for", e.Object, "with status", status, err.Error())
}

// Unpack json from request body
// Update column number in database and
// field num of json in database
func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	c := &models.Contact{}
	// Must valid new user Number.
	// User number must contains only numerics
	err := use.MatchNumber(e, c)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	models.UnpackingContact(c, e)
	err = database.UpdateNumber(e, c.Number)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(e.Object))
	log.Println("Number list updated general number by", e.Number, "new object", e.Object)
}

// this function unpacks json from request
// unpack JSON from this unpacked object
// and call target function from database package.
func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	c := &models.Contact{}
	err := use.MatchName(e, c)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	models.UnpackingContact(c, e)
	// call database package function using name of upgradable json field in database
	err = database.Update(e, "name", c.Name)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(e.Object))
	log.Println("Number list updated name by", e.Number, "new object", e.Object)

}

// this function updates general number of contct in db and
// same updates field num in json object in database
func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	c := &models.Contact{}
	models.UnpackingContact(c, e)
	err := database.Update(e, "nlist", c.NumberList)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	// write only modified field
	w.Write([]byte(e.Object))
	log.Println("Number list updated number list by", e.Number, "new object", e.Object)
}

// returns json string of contact from database(field object).
func InfoController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	err := database.GetInfo(e)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(e.Object))
	log.Println("Get info by", e.Number)
}

// Unpacks json from request body and delete target contact
// Write header bad request and error from database.Delete()
// if error isn't nil
func DeleteController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	err := database.Delete(e)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(e.Object))
	log.Println("Deleted", e.Object)
}

// Unpacks request body, call Create from database
// if Create() returns error, response includes this it
// and status code bad request
func CreateController(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		UnificatedErrorCapturing(e, w, e.Error, http.StatusBadRequest)
		return
	}
	err := use.MatchJsonFieldAndTarget(e)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	err = use.Matching(e)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	err = database.Create(e)
	if err != nil {
		UnificatedErrorCapturing(e, w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(e.Object))
	log.Println("Created", e.Object)
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
