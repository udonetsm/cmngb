package models

// This package can be imported as github.com/udonetsm/client/models.
// Server side uses it

import (
	"encoding/json"
)

// JSON object for making request to server side
// includes:
// target for fill entry_id in database
// Contact for build json string for use functions on the server side
type Entries struct {
	Number string `gorm:"number" json:"number"`
	// Object can be empty if using the DeleteOrInfo function.
	// See package github.com/udonetsm/client/http
	Object *Contact `gorm:"object" json:"object,omitempty"`
	// Error using on the server side for
	// answer about errors to clients...
	// default empty
	Error error `json:"error,omitempty" gorm:"-" `
}

// Pack object to json string
func (j *Entries) PackEntries(contact *Contact) (data []byte, err error) {
	if contact.Name == "" && contact.Number == "" && contact.NumberList == nil {
		contact = nil
	}
	j.Object = contact
	data, err = json.Marshal(j)
	if err != nil {
		j.Error = err
		return
	}
	return
}

func (j *Entries) UnpackEntries(data []byte) {
	err := json.Unmarshal(data, j)
	if err != nil {
		j.Error = err
		return
	}
}

type PackUnpackerEntries interface {
	UnpackEntries([]byte)
	PackEntries(*Contact) ([]byte, error)
}

func PackingEntries(pu PackUnpackerEntries, c *Contact) (data []byte, err error) {
	data, err = pu.PackEntries(c)
	return
}

func UnpackingEntries(pu PackUnpackerEntries, data []byte) {
	pu.UnpackEntries(data)
}

// Contact object
type Contact struct {
	Number     string   `json:"num,omitempty"`
	Name       string   `json:"name,omitempty"`
	NumberList []string `json:"nlist,omitempty"`
}

type PackUnpackerContact interface {
	UnpackContact(*Entries) []byte
}

func UnpackingContact(p PackUnpackerContact, e *Entries) {
	p.UnpackContact(e)
}

func (c *Contact) PackContact(e *Entries) (data []byte) {
	data, err := json.Marshal(e.Object)
	if err != nil {
		e.Error = err
	}
	return
}
