package database

import (
	"database/sql"
	"fmt"
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
	a := j.Object
	db, d := LoadDb()
	defer d.Close()
	db.Table("entries").Select("object").Where("number=?", j.Number).Scan(&j.Object)
	if a == j.Object {
		return ""
	}
	return j.Object
}

func Create(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	return db.Create(j).Error
}

func UpdateNumber(j *models.Entries, newvalue string) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Model(j).Where("number=?", j.Number).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", "num")+
			fmt.Sprintf(`'"%s"')`, newvalue))).UpdateColumn("number", newvalue)
	return tx.Error
}

func UpdateName(j *models.Entries, newvalue string) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Model(j).Where("number=?", j.Number).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", "name")+
			fmt.Sprintf(`'"%s"')`, newvalue)))
	return tx.Error
}

func Delete(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Where("number=?", j.Number).Delete(j)
	return tx.Error
}
