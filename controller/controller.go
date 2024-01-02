package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/server/database"
	"github.com/udonetsm/server/models"
	"github.com/udonetsm/server/use"
)

func Create(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	use.Match(e, use.ENME)
	use.Match(e, use.ENUM)
	use.Match(e, use.EQAL)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	database.Create(e)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	e.PackedObject = string(models.PackContact(e))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, e.PackedObject)
}

func Search(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	database.Search(e)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, e.ObjectList)
}

func Info(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	use.Match(e, use.ENUM)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	database.Info(e)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, e.PackedObject)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	database.Delete(e)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, e.PackedObject)
}

func UpdateGeneralNumber(w http.ResponseWriter, r *http.Request) {
	e := request(w, r)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	use.Match(e, use.ENUM)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	database.UpdateEntryNumber(e)
	if e.Error != nil {
		errs(w, e, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, e.PackedObject)
}

func UpdateName(w http.ResponseWriter, r *http.Request) {
}

func UpdateListNumber(w http.ResponseWriter, r *http.Request) {
}

func MW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Accept", "application/json")
		e := request(w, r)
		if e.Error != nil {
			errs(w, e, http.StatusBadRequest)
			return
		}
		use.Match(e, use.ENUM)
		if e.Error != nil {
			errs(w, e, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func errs(w http.ResponseWriter, e *models.Entries, status int) {
	w.WriteHeader(status)
	fmt.Fprintln(w, e.Error)
	log.Println(e.Error)
}

// This is a local function for unpacking jsons from requests
func request(w http.ResponseWriter, r *http.Request) *models.Entries {
	r.Header.Add("Accept", "appliction/json")
	r.Header.Add("Content-Type", "application/json")
	e := &models.Entries{}
	var req []byte
	req, e.Error = io.ReadAll(r.Body)
	if e.Error != nil {
		errs(w, e, http.StatusInternalServerError)
		return e
	}
	models.UnpackingEntries(e, req)
	r.Body = io.NopCloser(bytes.NewBuffer(req))
	return e
}
