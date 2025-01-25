package services

import (
	"gloriusaiapi/models"
	"github.com/jinzhu/gorm"
)

func CreateMessage(db *gorm.DB, userID uint, content, response string) (*models.Message, error) {
	message := models.Message{
		UserID:   userID,
		Content:  content,
		Response: response,
	}

	if err := db.Create(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func GetMessages(db *gorm.DB, userID uint) ([]models.Message, error) {
	var messages []models.Message
	if err := db.Where("user_id = ?", userID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
