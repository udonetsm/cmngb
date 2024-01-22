package models

import (
	"encoding/json"
	"io"
)

type Entries struct {
	// Id is a main field in the table. The field if primary key
	Id string `gorm:"id" json:"id,omitempty"`
	// Jcontact for packing contacts to the json object
	Jcontact *Contact `gorm:"-" json:"contact,omitempty"`
	// Contact for returning single jsons strings from the database
	Contact string `gorm:"contact" json:"-"`
	// Error here is for matching errors with nil and do something
	// if something went wrong somewhere.
	// But if somethig went wrong, should use Error.Error() as ErrMsg
	// for building json answer to the clients
	Error error `gorm:"-" json:"-"`
	// ErrMsg for building answer json to the clients
	ErrMsg string `gorm:"-" json:"error,omitempty"`
	// Contactlist is for collect all of found contacts in the Search function
	ContactList []string `gorm:"-" json:"contactlist,omitempty"`
	Owner       string   `gorm:"-" json:"-"`
	Secret      string   `gorm:"-" json:"-"`
	Token       string   `gorm:"-" json:"token,omitempty"`
	Ok          bool     `gorm:"-" json:"-"`
}

type Users struct {
	User_id   int `gorm:"-"`
	User_name string
	Secret    string
}

// This is a model of contact
type Contact struct {
	Name   string   `json:"name,omitempty"`
	Number string   `json:"number,omitempty"`
	List   []string `json:"list,omitempty"`
}

// Packs entry struct and writes packed entry to the some io.Writer
func (e *Entries) PackEntry(out io.Writer) {
	encoder := json.NewEncoder(out)
	encoder.Encode(e)
}

// Unpacks packedEntry to the entry struct
func (e *Entries) UnpackEntry(from []byte) {
	e.Error = json.Unmarshal(from, e)
}

// Packs some contact struct, writes error to the enctry error and retrun packed contact(json)
func (c *Contact) PackContact(e *Entries) []byte {
	data, err := json.Marshal(c)
	if err != nil {
		e.Error = err
	}
	return data
}

// Unpack contact object from packetdContact to the contact struct nd writes error to the entry error
func (contact *Contact) UnpackContact(packedContact []byte, entry *Entries) {
	entry.Error = json.Unmarshal(packedContact, contact)
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

// Use PackEntry method
func PackingEntry(entry PackUnpackerEntry, out io.Writer) {
	entry.PackEntry(out)
}

// Use UnpackEntry method
func UnpackingEntry(entry PackUnpackerEntry, packedEntry []byte) {
	entry.UnpackEntry(packedEntry)
}

// Use PackContact method
func PackingContact(contact PackUnpackerContact, e *Entries) (packedContact []byte) {
	packedContact = contact.PackContact(e)
	return
}

// Use UnpackContact method
func UnpackingContact(contact PackUnpackerContact, packedContact []byte, errTo *Entries) {
	contact.UnpackContact(packedContact, errTo)
}
