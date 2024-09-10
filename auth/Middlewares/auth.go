package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/metalpoch/olt-blueprint/auth/constants"
)

func TokenCheck(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.SALT), nil
	})
	return token, err
}
