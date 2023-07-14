package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/types"
	"net/http"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON("not authorized")
	}

	if !user.IsAdmin {
		return c.Status(http.StatusForbidden).JSON("forbidden")
	}

	return c.Next()
}
