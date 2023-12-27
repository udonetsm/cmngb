package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/udonetsm/client/models"
	"github.com/udonetsm/server/servermodels"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// This function returns object from database like JSON
func GetInfo(e *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Where("number = ?", e.Number).Find(e).Scan(e)
	if tx.RowsAffected < 1 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

// creates contact in database
func Create(e *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	return db.Create(e).Error
}

// This function updates general number phone indatabase
// And returns new object about contact from database if everything ok
func UpdateNumber(e *models.Entries, newvalue string) error {
	db, d := LoadDb()
	defer d.Close()
	// command UpdateColumn returns upgraded object from db
	// it can be pass by http to client
	tx := db.Model(e).Clauses(clause.Returning{Columns: []clause.Column{{Name: "object"}}}).
		Where("number=?", e.Number).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", "num")+
			fmt.Sprintf(`'"%s"')`, newvalue))).UpdateColumn("number", newvalue)
	if tx.RowsAffected == 0 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

// This function can upgrade all json field in database except number field
func Update(e *models.Entries, upgradableJSONfield string, newvalue interface{}) error {
	db, d := LoadDb()
	defer d.Close()
	// command Update returns upgraded object from db
	// it can be pass by http to client
	tx := db.Model(e).Where("number=?", e.Number).Clauses(clause.Returning{Columns: []clause.Column{{Name: "object"}}}).UpdateColumn("object",
		gorm.Expr("jsonb_set(object, "+fmt.Sprintf("'{%s}',", upgradableJSONfield)+
			fmt.Sprintf(`'"%s"')`, newvalue)))
	if tx.RowsAffected == 0 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}

// delete contact from database and return deleted object to pass it by http
func Delete(e *models.Entries) error {
	db, d := LoadDb()
	defer d.Close()
	tx := db.Where("number=?", e.Number).Find(e).Scan(e).Delete(e)
	if tx.RowsAffected < 1 {
		tx.Error = errors.New("NOT FOUND")
	}
	return tx.Error
}
