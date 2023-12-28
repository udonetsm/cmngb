package use

import (
	"errors"
	"regexp"

	"github.com/udonetsm/client/models"
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
		if ok := exp.Match([]byte(e.Object.Name)); !ok {
			e.Error = errors.New("INVALID NAME")
		}
	case "entrynumber":
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Number)); !ok {
			e.Error = errors.New("INVALID ENTRY NUMBER")
		}
	case "contactnumber":
		exp = regexp.MustCompile(regexpNumber)
		if ok := exp.Match([]byte(e.Object.Number)); !ok {
			e.Error = errors.New("CONTACT NUMBER INVALID")
		}
	case "contactandentrynumbers":
		if e.Number != e.Object.Number {
			e.Error = errors.New("GNUM AND CNUM AREN'T EQUAL")
		}
	}
}
