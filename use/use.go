package use

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/udonetsm/client/models"
)

const (
	reName   = "^[A-Z][a-z]* [A-Z][a-z]*$"
	reNumber = "^[0-9]*$"
)

func MatchName(e *models.Entries, c *models.Contact) (err error) {
	regName := regexp.MustCompile(reName)
	models.UnpackingContact(c, []byte(e.Object))
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
	models.UnpackingContact(c, []byte(e.Object))
	ok := regNumber.Match([]byte(c.Number))
	if !ok {
		err = errors.New("INVALID NUMBER")
	}
	return
}

// This function match internal Name and Number and returns error if one of them isn't valid
func Matching(j *models.Entries) (err error) {
	c := &models.Contact{}
	answer := ""
	var errs []error
	err = MatchName(j, c)
	if err != nil {
		errs = append(errs, err)
	}
	err = MatchNumber(j, c)
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
