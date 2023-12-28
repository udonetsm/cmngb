package database

import (
	"database/sql"
	"log"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/servermodels"

	"gorm.io/gorm"
)

// for load config for database connection
// configuration it and connect to database.
func LoadDb() (*gorm.DB, *sql.DB) {
	y := &servermodels.YAMLObject{}
	db := servermodels.LoadCfgAndGetDB(y, "/etc/cmngr/cfg.yaml")
	d, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	return db, d
}

func UpdateEntryNumber(e *models.Entries) (err error) {
	return
}

func UpdateEntryName(e *models.Entries) (err error) {
	return
}

func UpdateEntryListNumber(e *models.Entries) (err error) {
	return
}

func Info(e *models.Entries) (err error) {
	return
}

func Delete(e *models.Entries) (err error) {
	return
}

func Create(e *models.Entries) (err error) {
	return
}
