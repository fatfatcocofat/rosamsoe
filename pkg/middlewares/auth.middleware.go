package middlewares

import (
	"fmt"
	"strings"

	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/fatfatcocofat/rosamsoe/platform/database"
	"github.com/fatfatcocofat/rosamsoe/platform/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(config *config.Config, c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Failed when parsing auth token from request",
		})
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		logger.Debug().Err(err).Msg("Failed to parse jwt token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "The auth token provided has expired or is invalid",
		})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "The auth token provided has expired or is invalid",
		})
	}

	var user models.User
	database.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if user.ID.String() != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "The user belonging to this token no longger exists",
		})
	}

	c.Locals("user", models.UserFilterRecord(&user))

	return c.Next()
}

func UseAuthMiddleware(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return AuthMiddleware(config, c)
	}
}
