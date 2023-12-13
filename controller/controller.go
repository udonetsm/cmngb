package controller

import (
	"net/http"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/database"
)

func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	rj := &models.RequestJSON{}
	rj.UnpackRequest(r)
	database.UpdateNumber(rj)
}
