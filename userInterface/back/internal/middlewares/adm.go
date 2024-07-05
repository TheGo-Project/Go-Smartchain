package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/rosty-git/Smartchain-backend/internal/models"
)

// AdmMiddleware Middleware JWT function
func AdmMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		slog.Info("AdmMiddleware")

		user := c.Locals("User").(models.User)

		slog.Info("AdmMiddleware User", "user", user)

		if !user.Admin {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}
}
