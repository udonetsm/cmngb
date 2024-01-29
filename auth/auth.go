package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/udonetsm/cmngb/models"
)

func CreateToken(e *models.Entries, livetime int64) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: livetime,
		Id:        e.Owner,
	})
	e.Token, e.Error = unsignedToken.SignedString([]byte(e.Owner))
}

func TokenValid(e *models.Entries) bool {
	token, err := jwt.ParseWithClaims(e.Token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(e.Owner), nil
	})
	e.Error = err
	return token.Valid
}
