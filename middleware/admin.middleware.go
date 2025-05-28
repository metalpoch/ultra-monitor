package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func AdminAccess(ctx fiber.Ctx) error {
	if !ctx.Locals("is_admin").(bool) {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to access this resource."})
	}
	return ctx.Next()
}
