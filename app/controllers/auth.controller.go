package controllers

import (
	"strings"
	"time"

	"github.com/fatfatcocofat/rosamsoe/app/response"

	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/fatfatcocofat/rosamsoe/pkg/validator"
	"github.com/fatfatcocofat/rosamsoe/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Config *config.Config
	Logger *zerolog.Logger
}

func NewAuthController(config *config.Config, logger *zerolog.Logger) *AuthController {
	return &AuthController{
		Config: config,
		Logger: logger,
	}
}

// RegisterController handles user registration.
//
// @Summary Register a new user
// @Description Register a new user by providing user information in the request body.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.UserRegisterRequest true "User registration payload"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} response.BadRequest
// @Failure 500 {object} response.ServerError
// @Router /auth/register [post]
func (c *AuthController) RegisterController(ctx *fiber.Ctx) error {
	var payload *models.UserRegisterRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := validator.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Errors:  errors,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger.Error().Err(err).Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ServerError{
			Success: false,
			Message: response.SERVER_ERROR_MSG,
		})
	}

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
	}

	result := database.DB.Create(&newUser)

	if result.Error != nil && database.IsDuplicateError(result.Error) {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "User with that email already exists",
		})
	} else if result.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadRequest{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.Success{
		Success: true,
		Data:    models.UserFilterRecord(&newUser),
	})
}

// LoginController handles user login.
//
// @Summary Authenticate user and generate access token
// @Description Authenticate a user by providing login credentials in the request body and generate an access token.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.UserLoginRequest true "User login payload"
// @Success 200 {object} response.TokenResponse
// @Failure 400 {object} response.BadRequest
// @Failure 502 {object} response.BadGateway
// @Router /auth/login [post]
func (c *AuthController) LoginController(ctx *fiber.Ctx) error {
	var payload *models.UserLoginRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := validator.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Errors:  errors,
		})
	}

	var user models.User
	result := database.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": now.Add(c.Config.JwtExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(c.Config.JwtSecret))
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed generating jwt token")
		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.TokenResponse{
		Success: true,
		Data: response.TokenData{
			Token:     tokenStr,
			ExpiresIn: claims["exp"].(int64),
		},
	})
}
