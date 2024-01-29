package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/udonetsm/cmngb/auth"
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

func NewUser(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	database.NewUser(e)
	if e.Error != nil {
		fmt.Println(e.Error)
		errs(w, http.StatusBadRequest, e)
		return
	}
	auth.CreateToken(e, int64(time.Second*60*60*24*30))
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	models.PackingEntry(&models.Entries{Token: e.Token}, w)
	log.Printf("%s got token %s", e.Owner, e.Token)
}

// Get target id from request json and
// and get json of contact
// record from database using it
func Info(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	database.Info(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	ok(w, e, http.StatusFound)
}

// Get target name from request json and
// search all records in database using it
func Search(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	// call database function Search
	database.Search(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	e.Id = e.Jcontact.Name
	ok(w, e, http.StatusFound)
}

// Get target id from request json and
// delete record from database using it
func Delete(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	// call database function Delete
	database.Delete(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	ok(w, e, http.StatusOK)
}

// Get data from request json and
// add new record in database using it
func Create(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}

	use.Match(e, use.NAME)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	use.Match(e, use.EQAL)
	if e.Error != nil {
		e.Jcontact.Number = e.Id
		e.ErrMsg = e.Error.Error()
		e.Error = nil
	}
	e.Contact = string(models.PackingContact(e.Jcontact, e))
	// Pack all data to json
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	e.Jcontact = nil
	// Try to insert record in db
	database.Create(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	ok(w, e, http.StatusCreated)
}

// Updates target json field in database
// If name field in request json is empty
// updates list and vice versa
func Update(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	if len(e.Jcontact.Name) > 0 {
		use.Match(e, use.NAME)
		if e.Error != nil {
			errs(w, http.StatusBadRequest, e)
			return
		}
	}
	database.Update(e)
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	ok(w, e, http.StatusOK)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	e := &models.Entries{}
	request(w, r, e)
	e.Owner = r.Header.Get("owner")
	e.Secret = r.Header.Get("secret")
	auth.CreateToken(e, int64(time.Second*60*60*24*30))
	if e.Error != nil {
		errs(w, http.StatusBadRequest, e)
		return
	}
	models.PackingEntry(&models.Entries{Token: e.Token}, w)
	log.Printf("%s got token %s", e.Owner, e.Token)
}

func Valid(r *http.Request, e *models.Entries) {
	e.Token = r.Header.Get("Authorization")
	database.GetSecret(e)
	e.Ok = auth.TokenValid(e)
}

// Unpack json from request body, valid id and pass request to the next http.Handler
func MW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Accept", "application/json")
		e := &models.Entries{}
		request(w, r, e)
		if e.Error != nil {
			errs(w, http.StatusBadRequest, e)
			return
		}
		Valid(r, e)
		if !e.Ok {
			e.Error = errors.New("INVALID TOKEN")
			errs(w, http.StatusUnauthorized, e)
			return
		}
		if r.URL.RequestURI() != "/get/list" {
			use.Match(e, use.ENUM)
			if e.Error != nil {
				errs(w, http.StatusBadRequest, e)
				return
			}
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
		errs(w, http.StatusBadRequest, e)
		return
	}
	e.Ok = r.RequestURI == "/create/user"
	e.Owner = r.Header.Get("owner")
	models.UnpackingEntry(e, data)
	e.Contact = string(models.PackingContact(e.Jcontact, e))
	r.Body = io.NopCloser(bytes.NewBuffer(data))
}

// This is a local function which logging errors and
// writes it to the ResponseWriter...
func errs(w http.ResponseWriter, status int, e *models.Entries) {
	w.WriteHeader(status)
	// Logging errors and target objects from request
	log.Printf("[%v | ERROR: %v] ", e.Id, e.Error)
	// Pack and send json with some error to the client
	models.PackingEntry(&models.Entries{ErrMsg: e.Error.Error()}, w)
}

// This is a local function which logging success request/response-s
// and writes answer to the ResponseWriter
func ok(w http.ResponseWriter, e *models.Entries, code int) {
	// Logging success implementation with target as Id and error
	// Error logging in the create function if contact number and entry number aren't equal.
	w.WriteHeader(code)
	if len(e.ErrMsg) == 0 {
		e.ErrMsg = "no"
	}
	mes := fmt.Sprintf("OK for [%v] with %v errors", e.Id, e.ErrMsg)
	if len(e.ContactList) != 0 {
		mes += fmt.Sprintf(" and answer %v", e.ContactList)
		log.Println(mes)
		fmt.Fprintln(w, e.ContactList)
		return
	}
	mes += fmt.Sprintf(" and answer %v", e.Contact)
	log.Println(mes)
	fmt.Fprintln(w, e.Contact)
}
