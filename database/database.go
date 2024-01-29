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

func NewUser(e *models.Entries) {
	db, y := LoadDb(e)
	u := &models.Users{
		User_name: e.Owner,
		Secret:    e.Secret,
	}
	tx := db.Create(&u)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
	CreateTableForUser(db, e, y)
}

func LoadDb(e *models.Entries) (*gorm.DB, *YAMLObject) {
	y := &YAMLObject{}
	db := LoadCfgAndGetDB(y, "/etc/cfg.yaml", e)
	if y.Error != nil {
		e.Error = y.Error
		return new(gorm.DB), new(YAMLObject)
	}
	return db, y
}

func Info(e *models.Entries) {
	db, _ := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.Table(fmt.Sprintf("%s_entries", e.Owner)).First(&e)
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
	db, _ := LoadDb(e)
	if e.Error != nil {
		return
	}
	var rows *sql.Rows
	var err error
	if len(e.Jcontact.Name) == 0 {
		rows, err = db.Table(fmt.Sprintf("%s_entries", e.Owner)).Select("contact").Rows()
	} else {
		rows, err = db.Table(fmt.Sprintf("%s_entries", e.Owner)).
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
	db, _ := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.Table(fmt.Sprintf("%s_entries", e.Owner)).First(&e).Delete(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

// Search and return all contacts which contains
// target string in the name field. If Name field is
// empty, function returns all of contacts in the storage.
func Create(e *models.Entries) {
	db, _ := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.Table(fmt.Sprintf("%s_entries", e.Owner)).Create(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

func GetSecret(e *models.Entries) {
	db, _ := LoadDb(e)
	if e.Error != nil {
		return
	}
	tx := db.Table("users").Select("secret").Where("user_name=?", e.Owner).Scan(&e.Secret)
	e.Error = tx.Error
}

// Make command update for call gorm function Exec or Raw.
func buildCmd(e *models.Entries) (cmd bytes.Buffer) {
	cmd.WriteString(fmt.Sprintf("update %s_entries set ", e.Owner))
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
	db, _ := LoadDb(e)
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
