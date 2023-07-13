package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/types"
)

func getAuthUser(c *fiber.Ctx) (*types.User, bool) {
	user, ok := c.Context().UserValue("user").(*types.User)
	return user, ok
}
