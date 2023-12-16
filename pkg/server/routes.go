package server

import (
	"fmt"

	"github.com/fatfatcocofat/rosamsoe/app/controllers"
	"github.com/fatfatcocofat/rosamsoe/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *Server) routes() {
	s.App.Get("/", func(c *fiber.Ctx) error {
		swaggerURL := fmt.Sprintf("http://%s:%d/swagger", s.Config.ServerHost, s.Config.ServerPort)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Welcome to Rosamsoe, Developed With 3> By Fathurrohman.",
			"repository": fiber.Map{
				"type": "git",
				"url":  "https://github.com/fatfatcocofat/rosamsoe",
			},
			"documentation": fiber.Map{
				"name": "Rosamsoe API",
				"url":  swaggerURL,
			},
		})
	})

	v1 := s.App.Group("/api/v1")

	authController := controllers.NewAuthController(s.Config, s.Logger)
	v1.Route("/auth", func(router fiber.Router) {
		router.Post("/register", authController.RegisterController)
		router.Post("/login", authController.LoginController)
	})

	userController := controllers.NewUserController(s.Config, s.Logger)
	v1.Route("/user", func(router fiber.Router) {
		router.Use(middlewares.UseAuthMiddleware(s.Config))
		router.Get("/", userController.InfoController)
	})

	walletController := controllers.NewWalletController(s.Config, s.Logger)
	v1.Route("/wallet", func(router fiber.Router) {
		router.Use(middlewares.UseAuthMiddleware(s.Config))
		router.Get("/", walletController.ListController)
		router.Post("/", walletController.CreateController)

		router.Get("/:address", walletController.ShowController)
		router.Delete("/:address", walletController.DeleteController)
		router.Patch("/:address", walletController.UpdateController)
	})
}

func (s *Server) catchNotFoundRoutes() {
	s.App.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Endpoint not found",
		})
	})
}

func (s *Server) handleSwaggerRoutes() {
	s.App.Get("/swagger/*", swagger.HandlerDefault)
}
