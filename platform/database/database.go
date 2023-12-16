package database

import (
	"fmt"

	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	log "github.com/fatfatcocofat/rosamsoe/platform/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *config.Config) (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Info().Msg("Running Migrations")

	err = DB.AutoMigrate(&models.User{}, &models.Wallet{})
	if err != nil {
		log.Fatal().Err(err).Msg("Migration Failed")
	}

	log.Info().Msg("Connected Successfully to the Database")

	return
}
