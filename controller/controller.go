package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
)

func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	a := &models.RequestJSON{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, a)
	if err != nil {
		log.Fatal(err)
	}
	database.Unumber()
}
