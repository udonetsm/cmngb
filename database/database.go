package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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
		log.Println(y.Error)
		e.Error = gorm.ErrInvalidDB
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

func DeleteById(e *models.Entries) {
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

func makeList(list []string) (value bytes.Buffer) {
	value.WriteString("[")
	for i, v := range list {
		value.WriteString("\"" + v + "\"")
		if i == len(list)-1 {
			value.WriteString("]")
			break
		}
		value.WriteString(",")
	}
	return
}

func buildCommandArray(e *models.Entries, u string) (cmd bytes.Buffer) {
	cmd.WriteString("update entries set contact=")
	if u == UNLST {
		list := makeList(e.Jcontact.List)
		cmd.WriteString(fmt.Sprintf("jsonb_set(contact, '{%s}', '%v')", UNLST, list.String()))
	} else if u == UNAME {
		cmd.WriteString(fmt.Sprintf("jsonb_set(contact, '{%s}', '\"%s\"')", UNAME, e.Jcontact.Name))
	} else if u == UNUMB {
		cmd.WriteString(fmt.Sprintf("jsonb_set(contact, '{%s}', '%s'), id='%s'", UNUMB, e.Jcontact.Number, e.Jcontact.Number))
	}
	cmd.WriteString(fmt.Sprintf("where id='%s' returning contact", e.Id))
	return
}

func Update(e *models.Entries, u string) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	cmd := buildCommandArray(e, u)
	if len(cmd.String()) == 0 {
		e.Error = gorm.ErrInvalidField
		return
	}
	rows, err := db.Raw(cmd.String()).Rows()
	if err != nil {
		e.Error = err
		return
	}
	if rows.Next() {
		err := rows.Scan(&e.Contact)
		if err != nil {
			e.Error = err
		}
		return
	}
	e.Error = gorm.ErrRecordNotFound
}
