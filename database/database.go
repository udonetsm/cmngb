package database

import (
	"fmt"

	"github.com/udonetsm/cmngb/models"
	"gorm.io/gorm"
)

func LoadDb() *gorm.DB {
	y := &YAMLObject{}
	db := LoadCfgAndGetDB(y, "/etc/cmngr/cfg.yaml")
	return db
}

func Info(e *models.Entries) {
	db := LoadDb()
	tx := db.First(e, "id=?", e.Id)
	e.Error = tx.Error
}

func Search(e *models.Entries) {
	db := LoadDb()
	rows, err := db.Model(&models.Entries{}).
		Select("contact").
		Where("contact->>'name' like ?", "%"+e.Jcontact.Name+"%").
		Rows()
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
	db := LoadDb()
	tx := db.First(&e).Delete(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

func Create(e *models.Entries) {
	db := LoadDb()
	tx := db.Create(&e)
	if tx.Error != nil {
		e.Error = tx.Error
		return
	}
}

func UpdateNumber(e *models.Entries) {
	db := LoadDb()
	rows, err := db.Raw("update entries set contact=" +
		fmt.Sprintf("jsonb_set(contact, '{number}', '%v'), id='%s' where id='%s' returning contact",
			e.Jcontact.Number, e.Jcontact.Number, e.Id)).Rows()
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

func Update(e *models.Entries, u string, save any) {
	db := LoadDb()
	cmd := fmt.Sprintf("update entries set contact=jsonb_set(contact, '{%s}', '\"%v\"') where id='%v' returning contact",
		u, save, e.Id)
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
