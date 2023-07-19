package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/types"
)

type resourceResponse struct {
	Results int   `json:"results"`
	Page    int64 `json:"page"`
	Data    any   `json:"data"`
}

func getAuthUser(c *fiber.Ctx) (*types.User, bool) {
	user, ok := c.Context().UserValue("user").(*types.User)
	return user, ok
}
