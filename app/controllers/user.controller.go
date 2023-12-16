package controllers

import (
	"github.com/fatfatcocofat/rosamsoe/app/response"

	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type UserController struct {
	Config *config.Config
	Logger *zerolog.Logger
}

func NewUserController(config *config.Config, logger *zerolog.Logger) *UserController {
	return &UserController{
		Config: config,
		Logger: logger,
	}
}

// InfoController retrieves user information.
//
// @Summary Get user information
// @Description Retrieve information about the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} response.Unauthorized
// @Failure 500 {object} response.ServerError
// @Router /auth/user [get]
func (c *UserController) InfoController(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.UserResponse)
	return ctx.JSON(response.Success{
		Success: true,
		Data: fiber.Map{
			"user": user,
		},
	})
}
