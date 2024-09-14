package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/metalpoch/olt-blueprint/auth/constants"
)

func CreateJWT(secret []byte, id uint, is_admin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["is_admin"] = is_admin
	claims["exp"] = time.Now().Add(constants.EXPIRE_JWT).Unix()

	return token.SignedString(secret)
}
