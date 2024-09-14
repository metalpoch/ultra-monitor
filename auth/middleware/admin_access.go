package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/constants"
)

func AdminAccess(ctx fiber.Ctx) error {
	if !ctx.Locals("is_admin").(bool) {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"error": constants.FORBIDDEN_RESPONSE})
	}
	return ctx.Next()
}
