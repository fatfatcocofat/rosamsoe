package utils

import (
	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/gofiber/fiber/v2"
)

func ParseUserFromCtx(c *fiber.Ctx) models.UserResponse {
	user := c.Locals("user").(models.UserResponse)

	return user
}
