package database

import (
	"database/sql"
	"log"

	"github.com/udonetsm/client/models"
	"gorm.io/gorm"
)

func LoadDb() (*gorm.DB, *sql.DB) {
	y := &models.YAMLObject{}
	db := models.LoadCfgAndGetDB(y, "database/cfg.yaml")
	d, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	return db, d
}

func GetInfo(j *models.Entries) string {
	db, d := LoadDb()
	defer d.Close()
	db.Table("entries").Select("object").Where("number=?", j.Number).Scan(&j.Object)
	return j.Object
}

func Create(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	return db.Create(j).Error
}

func UpdateNumber(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	return db.Update("", "").Error
}
