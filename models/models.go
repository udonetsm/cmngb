package models

import (
	"encoding/json"
	"errors"
)

type Entries struct {
	Id          string   `gorm:"id" json:"id,omitempty"`
	Jcontact    *Contact `gorm:"-" json:"contact,omitempty"`
	Contact     string   `gorm:"contact" json:"-"`
	Error       error    `gorm:"-" json:"error,omitempty"`
	ContactList []string `gorm:"-" json:"contactlist,omitempty"`
}

type Contact struct {
	Name   string   `json:"name,omitempty"`
	Number string   `json:"number,omitempty"`
	List   []string `json:"list,omitempty"`
}

type PackUnpackerEntry interface {
	PackEntrie(c *Contact) []byte
	UnpackEntry(data []byte)
}

func (e *Entries) PackEntrie(c *Contact) (data []byte) {
	data, err := json.Marshal(e.Jcontact)
	if err != nil {
		e.Error = err
		return nil
	}
	return
}

func PackingEntry(pu PackUnpackerEntry, c *Contact) (data []byte) {
	data = pu.PackEntrie(c)
	return
}

func (e *Entries) UnpackEntry(data []byte) {
	err := json.Unmarshal(data, e)
	if err != nil {
		e.Error = errors.New("INVALID JSON")
		return
	}
	e.Error = err
}

func UnpackingEntry(pu PackUnpackerEntry, data []byte) {
	pu.UnpackEntry(data)
}
