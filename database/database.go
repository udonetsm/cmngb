package database

import (
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

func makeCmd(e *models.Entries, u string) string {
	switch u {
	case UNAME:
		return fmt.Sprintf("update entries set contact=jsonb_set(contact, '{%s}', '%s') where id='%s' returning contact",
			u, e.Jcontact.Name, e.Id)
	case UNUMB:
		return fmt.Sprintf("update entries set contact=jsonb_set(contact, '{%s}', '%v'), id='%s' where id='%s' returning contact",
			u, e.Jcontact.Number, e.Jcontact.Number, e.Id)
	case UNLST:
		s := "["
		for i, v := range e.Jcontact.List {
			if i == len(e.Jcontact.List)-1 {
				s += "\"" + v + "\"]"
				return fmt.Sprintf("update entries set contact=jsonb_set(contact, '{%s}', '%v') where id='%s' returning contact",
					u, s, e.Id)
			}
			s += "\"" + v + "\" ,"
		}
	}
	e.Error = gorm.ErrRecordNotFound
	return ""
}

func Update(e *models.Entries, u string) {
	db := LoadDb(e)
	if e.Error != nil {
		return
	}
	cmd := makeCmd(e, u)
	rows, err := db.Raw(cmd).Rows()
	if err != nil {
		e.Error = err
		return
	}
	if rows.Next() {
		err := rows.Scan(&e.Contact)
		if err != nil {
			e.Error = err
			return
		}
	}
}
