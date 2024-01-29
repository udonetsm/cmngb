package use

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/udonetsm/cmngb/models"
)

// For use it in regexp function
const (
	regexpName   = "^[A-Z][a-z]* [A-Z][a-z]*"
	regexpNumber = "^[0-9]*$"
	ENUM         = "entrynumber"
	NAME         = "name"
	CNUM         = "contactnumber"
	EQAL         = "contactandentrynumbers"
)

func Match(e *models.Entries, matchable string) {
	var exp *regexp.Regexp
	switch matchable {
	case NAME:
		exp = regexp.MustCompile(regexpName)
		if ok := exp.Match([]byte(e.Jcontact.Name)); !ok {
			e.Error = errors.New("INVALID NAME")
			return
		}
	case ENUM:
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Id)); !ok || e.Id == "" {
			e.Error = errors.New("INVALID ENTRY NUMBER")
			return
		}
	case CNUM:
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Jcontact.Number)); !ok {
			e.Error = errors.New("INVALID CONTACT NUMBER")
			return
		}
	case EQAL:
		if e.Id != e.Jcontact.Number {
			e.Error = fmt.Errorf("%s %s AND %s %s AREN'T EQUAL(entry number has been set like contact number)", ENUM, e.Id, CNUM, e.Jcontact.Number)
			return
		}
	}
}
