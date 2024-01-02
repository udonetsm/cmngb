package use

import (
	"errors"
	"regexp"

	"github.com/udonetsm/cmngb/models"
)

// For use it in regexp function
const (
	regexpName   = "^[A-Z][a-z]* [A-Z][a-z]*$"
	regexpNumber = "^[0-9]*$"
	ENUM         = "entrynumber"
	ENME         = "name"
	CNUM         = "contactnumber"
	EQAL         = "contactandentrynumbers"
)

func Match(e *models.Entries, matchable string) {
	var exp *regexp.Regexp
	switch matchable {
	case "name":
		exp = regexp.MustCompile(regexpName)
		if ok := exp.Match([]byte(e.Jcontact.Name)); !ok {
			e.Error = errors.New("INVALID NAME")
		}
	case "entrynumber":
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Id)); !ok {
			e.Error = errors.New("INVALID ENTRY NUMBER")
		}
	case "contactnumber":
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Jcontact.Number)); !ok {
			e.Error = errors.New("INVALID CONTACT NUMBER")
		}
	case "contactandentrynumbers":
		if e.Id != e.Jcontact.Number {
			e.Error = errors.New(ENUM + e.Id + " AND " + CNUM + e.Jcontact.Number + " AREN'T EQUAL")
		}
	}
}
