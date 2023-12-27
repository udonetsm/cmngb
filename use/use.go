package use

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/udonetsm/client/models"
)

// Function for matching internal target number
// and JSON object number. If aren't equal write error in ResponseWriter
func MatchJsonFieldAndTarget(e *models.Entries) (err error) {
	c := &models.Contact{}
	models.UnpackingContact(c, e)
	if c.Number != e.Number {
		err = errors.New("TARGET NUMBER AND JSON OBJECT NUMBER AREN'T EQUAL")
	}
	return
}

// For use it in regexp function
const (
	reName   = "^[A-Z][a-z]* [A-Z][a-z]*$"
	reNumber = "^[0-9]*$"
)

// For matching name using regexp.
// Name must be like "Name Surname" only
func MatchName(e *models.Entries, c *models.Contact) (err error) {
	regName := regexp.MustCompile(reName)
	models.UnpackingContact(c, e)
	ok := regName.Match([]byte(c.Name))
	if !ok {
		err = errors.New("INVALID NAME")
	}
	return
}

// Mathc number only and returns error if it isn't valid
func MatchNumber(e *models.Entries, c *models.Contact) (err error) {
	regNumber := regexp.MustCompile(reNumber)
	// Unpack Object fiels of models.Entries to contact
	models.UnpackingContact(c, e)
	ok := regNumber.Match([]byte(c.Number))
	if !ok {
		err = errors.New("INVALID NUMBER")
	}
	return
}

// This function match internal Name and Number and returns error if one of them isn't valid
func Matching(e *models.Entries) (err error) {
	c := &models.Contact{}
	answer := ""
	var errs []error
	err = MatchName(e, c)
	if err != nil {
		errs = append(errs, err)
	}
	err = MatchNumber(e, c)
	if err != nil {
		errs = append(errs, err)
	}
	if errs != nil {
		for _, v := range errs {
			// build answer with all of errors
			answer += fmt.Sprintf("[%v] ", v)
		}
		// build error with answer
		err = errors.New(answer)
	}
	return
}
