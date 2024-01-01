package database

import (
	"database/sql"
	"log"

	"github.com/udonetsm/server/cfgsrv"
	"github.com/udonetsm/server/models"

	"gorm.io/gorm"
)

// for load config for database connection
// configuration it and connect to database.
func LoadDb() (*gorm.DB, *sql.DB) {
	y := &cfgsrv.YAMLObject{}
	db := cfgsrv.LoadCfgAndGetDB(y, "/etc/cmngr/cfg.yaml")
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

func Create(e *models.Entries) {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Exec("insert into entries(number, object) values(?,?)", e.Number, e.Object)
	e.Error = tx.Error
}

func Info(e *models.Entries) {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Table("entries").
		Select("object").
		Where("number=?", e.Number).
		Scan(e.PackedObject)
	e.Error = tx.Error
}

func Search(e *models.Entries) {
	db, d := LoadDb()
	defer d.Close()
	rows, err := db.Table("entries").
		Select("object").
		Where("object->>'name' like ?", "%"+e.Object.Name+"%").
		Rows()
	e.Error = err
	for rows.Next() {
		a := ""
		rows.Scan(&a)
		e.ObjectList = append(e.ObjectList, "\n", a)
	}
}

func Delete(e *models.Entries) {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Select("object").
		Table("entries").
		Where("number=?", e.Number).
		Scan(&e.PackedObject).
		Delete(e)
	e.Error = tx.Error
}
