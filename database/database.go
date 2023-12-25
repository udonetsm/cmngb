package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/udonetsm/client/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LoadDb() (*gorm.DB, *sql.DB) {
	y := &models.YAMLObject{}
	db := models.LoadCfgAndGetDB(y, "/etc/cmngr/cfg.yaml")
	d, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	return db, d
}

func GetInfo(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Where("number = ?", j.Number).Find(j).Scan(j)
	if tx.RowsAffected < 1 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

// creates contact in database
func Create(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	return db.Create(j).Error
}

func UpdateNumber(j *models.Entries, newvalue string) error {
	db, d := LoadDb()
	defer d.Close()
	// command UpdateColumn returns upgraded object from db
	// it can be pass by http to client
	tx := db.Model(j).Clauses(clause.Returning{Columns: []clause.Column{{Name: "object"}}}).
		Where("number=?", j.Number).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", "num")+
			fmt.Sprintf(`'"%s"')`, newvalue))).UpdateColumn("number", newvalue)
	if tx.RowsAffected == 0 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

func Update(j *models.Entries, upgradableJSONfield string, newvalue interface{}) error {
	db, d := LoadDb()
	defer d.Close()
	// command Update returns upgraded object from db
	// it can be pass by http to client
	tx := db.Model(j).Where("number=?", j.Number).Clauses(clause.Returning{Columns: []clause.Column{{Name: "object"}}}).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", upgradableJSONfield)+
			fmt.Sprintf(`'"%s"')`, newvalue)))
	if tx.RowsAffected == 0 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

// delete contact from database and return deleted object to pass it by http
func Delete(j *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Where("number=?", j.Number).Find(j).Scan(j).Delete(j)
	if tx.RowsAffected < 1 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}
