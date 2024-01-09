package models

import (
	"encoding/json"
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

type EntryUnpacker interface {
	UnpackEntry(data []byte)
}

type EntryPacker interface {
	PackEntrie(c *Contact) []byte
}

func (e *Entries) PackEntrie(c *Contact) (data []byte) {
	data, err := json.Marshal(e.Jcontact)
	if err != nil {
		e.Error = err
		return nil
	}
	return
}

func PackingEntry(pu EntryPackUnpacker, c *Contact) (data []byte) {
	data = pu.PackEntrie(c)
	return
}

func (e *Entries) UnpackEntry(data []byte) {
	err := json.Unmarshal(data, e)
	if err != nil {
		e.Error = err
		return
	}
	e.Error = err
}

func EntryUnpacking(pu EntryPackUnpacker, data []byte) {
	pu.UnpackEntry(data)
}

type EntryPackUnpacker interface {
	EntryPacker
	EntryUnpacker
}
