package services

import (
	"github.com/jinzhu/gorm"
	"gloriusaiapi/models"
)

func CreateOrUpdateContext(db *gorm.DB, userID uint, key, value string) (*models.Context, error) {
	var context models.Context

	err := db.Where("user_id = ? AND key = ?", userID, key).First(&context).Error
	if err == nil {
		context.Value = value
		if err := db.Save(&context).Error; err != nil {
			return nil, err
		}
	} else {
		context = models.Context{
			UserID: userID,
			Key:    key,
			Value:  value,
		}
		if err := db.Create(&context).Error; err != nil {
			return nil, err
		}
	}

	return &context, nil
}

func GetContexts(db *gorm.DB, userID uint) ([]models.Context, error) {
	var contexts []models.Context
	if err := db.Where("user_id = ?", userID).Find(&contexts).Error; err != nil {
		return nil, err
	}
	return contexts, nil
}
