package repository

import (
    "gloriusaiapi/config"
    "gloriusaiapi/models"
)

func CreateMessage(userID uint, content, response string) (*models.Message, error) {
    message := models.Message{
        UserID:   userID,
        Content:  content,
        Response: response,
    }

    if err := config.DB.Create(&message).Error; err != nil {
        return nil, err
    }
    return &message, nil
}

func GetMessagesByUser(userID uint) ([]models.Message, error) {
    var messages []models.Message
    if err := config.DB.Where("user_id = ?", userID).Find(&messages).Error; err != nil {
        return nil, err
    }
    return messages, nil
}
