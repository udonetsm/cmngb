package models

import (
	"encoding/json"
	"io"
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

func (e *Entries) PackEntry(out io.Writer) {
	encoder := json.NewEncoder(out)
	e.Error = encoder.Encode(e)
}

func (e *Entries) UnpackEntry(data []byte) {
	e.Error = json.Unmarshal(data, e)
}

func (c *Contact) PackContact(e *Entries) []byte {
	data, err := json.Marshal(c)
	e.Error = err
	return data
}

func (c *Contact) UnpackContact(data []byte, e *Entries) {
	e.Error = json.Unmarshal(data, c)
	if e.Error != nil {
		return
	}
}

type ContactPacker interface {
	PackContact(*Entries) []byte
}

type ContactUnpacker interface {
	UnpackContact([]byte, *Entries)
}

type EntryPacker interface {
	PackEntry(io.Writer)
}

type EntryUnpacker interface {
	UnpackEntry([]byte)
}

type PackUnpackerContact interface {
	ContactPacker
	ContactUnpacker
}

type PackUnpackerEntry interface {
	EntryPacker
	EntryUnpacker
}

func PakcingEntry(pue PackUnpackerEntry, out io.Writer) {
	pue.PackEntry(out)
}
func UnpackingEntry(pue PackUnpackerEntry, data []byte) {
	pue.UnpackEntry(data)
}
func PackingContact(puc PackUnpackerContact, e *Entries) (data []byte) {
	data = puc.PackContact(e)
	return
}
func UnpackingContact(puc PackUnpackerContact, data []byte, e *Entries) {
	puc.UnpackContact(data, e)
}
