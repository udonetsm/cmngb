package database

import (
	"github.com/udonetsm/client/models"
	"gorm.io/gorm"
)

func LoadDb() *gorm.DB {
	y := &models.YAMLObject{}
	db := models.LoadCfgAndGetDB(y, "database/cfg.yaml")
	return db
}

func Create(j *models.JSONObject) error {
	db := LoadDb()
	return db.Create(j).Error
}
