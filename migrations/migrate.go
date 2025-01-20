package migrations

import (
	"gloriusaiapi/config"
	"gloriusaiapi/models"
)

func Migrate() {
	config.DB.AutoMigrate(&models.User{}, &models.Message{}, &models.Context{})
}
