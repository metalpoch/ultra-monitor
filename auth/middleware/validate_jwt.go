package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(secret []byte) fiber.Handler {

	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"error": "authorization token invalid"})
		}
		receivedToken := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		ctx.Locals("id", claims["id"])
		ctx.Locals("is_admin", claims["is_admin"])

		return ctx.Next()
	}
}
