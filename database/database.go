package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/udonetsm/cmngb/models"
	"gorm.io/gorm"
)

const (
	UNAME = "name"
	UNUMB = "number"
	UNLST = "list"
)

func LoadDb(e *models.Entries) *gorm.DB {
	y := &YAMLObject{}
	db := LoadCfgAndGetDB(y, "/etc/cfg.yaml")
	if y.Error != nil {
		e.Error = y.Error
		return new(gorm.DB)
	}
	return db
}

func Info(e *models.Entries) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.First(e)
	e.Error = tx.Error
}

func JSONvalidator(data []byte, a any) bool {
	err := json.Unmarshal(data, a)
	return err == nil
}

// Search and return all contacts which contains
// target string in the name field. If Name field is
// empty, function returns all of contacts in the storage.
func Search(e *models.Entries) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	var rows *sql.Rows
	var err error
	if len(e.Jcontact.Name) == 0 {
		rows, err = db.Model(&models.Entries{}).Select("contact").Rows()
	} else {
		rows, err = db.Model(&models.Entries{}).
			Select("contact").
			Where("contact->>'name' like ?", "%"+e.Jcontact.Name+"%").
			Rows()
	}
	if err != nil {
		e.Error = err
		return
	}
	for rows.Next() {
		rows.Scan(&e.Contact)
		e.ContactList = append(e.ContactList, e.Contact)
	}
}

func Delete(e *models.Entries) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.First(&e).Delete(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

// Search and return all contacts which contains
// target string in the name field. If Name field is
// empty, function returns all of contacts in the storage.
func Create(e *models.Entries) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.Create(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

// Make command update for call gorm function Exec or Raw.
func buildCmd(e *models.Entries) (cmd bytes.Buffer) {
	cmd.WriteString("update entries set ")
	if len(e.Jcontact.Number) == 0 {
		cmd.WriteString(fmt.Sprintf("id='%s', ", e.Id))
	} else {
		cmd.WriteString(fmt.Sprintf("id='%s', ", e.Jcontact.Number))
	}
	cmd.WriteString(fmt.Sprintf("contact=contact||'%s'", e.Contact))
	cmd.WriteString(fmt.Sprintf(" where id='%s' ", e.Id))
	cmd.WriteString("returning contact")
	return
}

// Updates json object in db. Updates only fields got from request. Other fields doesn't update
func Update(e *models.Entries) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	cmd := buildCmd(e)
	rows, err := db.Raw(cmd.String()).Rows()
	if err != nil {
		e.Error = err
		return
	}
	e.Contact = ""
	if rows.Next() {
		e.Error = rows.Scan(&e.Contact)
		return
	} else {
		e.Error = gorm.ErrRecordNotFound
	}
}
