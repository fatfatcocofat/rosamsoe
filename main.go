package main

import (
	_ "github.com/fatfatcocofat/rosamsoe/docs"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/fatfatcocofat/rosamsoe/pkg/server"
	"github.com/fatfatcocofat/rosamsoe/platform/database"
	"github.com/fatfatcocofat/rosamsoe/platform/logger"
)

// @title Rosamsoe API
// @version 1.0
// @description Rosamsoe API Documentation
// @contact.name Fathurrohman
// @contact.url https://t.me/fatfatcocofat
// @license.name MIT License
// @license.url https://github.com/fatfatcocofat/rosamsoe/blob/main/LICENSE
// @host localhost:8000
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load environment variables")
	}

	err = database.ConnectDB(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to the Database")
	}

	app := server.New(&cfg)
	app.Serve()
}
