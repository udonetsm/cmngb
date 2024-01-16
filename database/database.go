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

func buildCommandArray(e *models.Entries, u string) (cmd []string) {
	cmd = append(cmd, "update entries set contact=contact || '{")
	if len(e.Jcontact.List) > 0 {
		list := makeList(e.Jcontact.List)
		cmd = append(cmd, fmt.Sprintf("\"%s\": %v", UNLST, list.String()))
	}
	if len(e.Jcontact.Name) > 0 {
		cmd = append(cmd, fmt.Sprintf("\"%s\": \"%s\"", UNAME, e.Jcontact.Name))
	}
	if len(e.Jcontact.Number) > 0 {
		cmd = append(cmd, fmt.Sprintf("\"%s\": \"%s\"", UNUMB, e.Jcontact.Number))
		cmd = append(cmd, fmt.Sprintf("}', id='%s' where id='%s' returnning contact",
			e.Jcontact.Number, e.Jcontact.Number))
		return
	}
	cmd = append(cmd, fmt.Sprintf("}' where id='%s' returning contact", e.Id))
	return
}

func Update(e *models.Entries, u string) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	cmd := buildCommandArray(e, u)
	if len(cmd) == 0 {
		e.Error = gorm.ErrInvalidField
		return
	}
	fmt.Println(cmd)
	return
	rows, err := db.Raw("").Rows()
	if err != nil {
		e.Error = err
		return
	}
	if rows.Next() {
		e.Error = rows.Scan(&e.Contact)
		return
	}
	e.Error = gorm.ErrRecordNotFound
}

func u(e *models.Entries) {

}
