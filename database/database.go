package database

import (
	"log"

	"github.com/udonetsm/client/models"
)

func UpdateNumber(rj *models.RequestJSON) {
	y := &models.YAMLObject{}
	y.YAMLCfg("database/cfg.yaml")
	log.Println(y, rj)
}
