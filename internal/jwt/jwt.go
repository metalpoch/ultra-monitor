package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, id int32, is_admin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(4 * time.Hour).Unix()
	return token.SignedString(secret)
}
