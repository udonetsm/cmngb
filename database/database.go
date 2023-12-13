package database

import (
	"log"

	"github.com/udonetsm/client/models"
)

func UpdateNumber(rj *models.RequestJSON) {
	y := &models.YAMLObject{}
	db := models.LoadCfgAndGetDB(y, "database/cfg.yaml")
	log.Println(db)
}
