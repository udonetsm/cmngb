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

func GetInfo(j *models.Entries) string {
	db := LoadDb()
	db.Table("entries").Select("object").Where("number=?", j.Number).Scan(&j.Object)
	return j.Object
}

func Create(j *models.Entries) error {
	db := LoadDb()
	return db.Create(j).Error
}

func UpdateNumber(j *models.Entries) error {
	db := LoadDb()
	return db.Update("", "").Error
}
