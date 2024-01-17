package controllers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/cmngb/database"
	"github.com/udonetsm/cmngb/models"
	"github.com/udonetsm/cmngb/use"
)

const (
	// Index of contact in arguments array
	CONTACT = iota
	// Index of target in arguments array
	TARGET
	// Index of error in arguments array
	ERROR
)

// Get target id from request json and
// and get json of contact
// record from database using it
func Info(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	database.Info(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	ok(w, e.Contact, e.Id)
}

// Get target name from request json and
// search all records in database using it
func Search(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Contact, e.Error)
		return
	}
	searchBy := e.Jcontact.Name
	// call database function Search
	database.Search(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Contact, e.Error)
		return
	}
	ok(w, e.ContactList, searchBy)
}

// Get target id from request json and
// delete record from database using it
func DeleteById(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	// call database function Delete
	database.DeleteById(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	ok(w, e.Contact, e.Id)
}

// Get data from request json and
// add new record in database using it
func Create(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}

	use.Match(e, use.ENME)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	use.Match(e, use.EQAL)
	if e.Error != nil {
		fmt.Fprintln(w, e.Error)
		e.Jcontact.Number = e.Id
	}
	e.Contact = string(models.PackingContact(e.Jcontact, e))
	// Pack all data to json
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	e.Jcontact = nil
	// Try to insert record in db
	database.Create(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	ok(w, e.Contact, e.Id)
}

// Updates target json field in database
// If name field in request json is empty
// updates list and vice versa
func Update(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	database.Update(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	ok(w, e.Contact, e.Id)
}

func MW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Accept", "application/json")
		e := &models.Entries{}
		request(w, r, e)
		if e.Error != nil {
			errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
			return
		}
		use.Match(e, use.ENUM)
		if e.Error != nil {
			errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// This is a local funtion for get data from request.
// If something went wrong it writes error in models.Entries{} object
// logging this object and writes it to the ResponseWriter
func request(w http.ResponseWriter, r *http.Request, e *models.Entries) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		errs(w, http.StatusBadRequest, nil, e.Id, e.Error)
		return
	}
	models.UnpackingEntry(e, data)
	cdata := models.PackingContact(e.Jcontact, e)
	e.Contact = string(cdata)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
}

// This is a local function which logging errors and
// writes it to the ResponseWriter...
func errs(w http.ResponseWriter, status int, a ...any) {
	fmt.Fprintln(w, a[ERROR])
	log.Printf("[%v | %v] ", a[ERROR], a[TARGET])
}

// This is a local function which logging success request/response-s
// and writes answer to the ResponseWriter
func ok(w http.ResponseWriter, a ...any) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, a[CONTACT])
	log.Printf("[OK FOR TARGET %v] %v", a[TARGET], a[CONTACT])
}
