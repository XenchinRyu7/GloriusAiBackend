package repository

import (
    "gloriusaiapi/config"
    "gloriusaiapi/models"
)

func CreateUser(username, password string) (*models.User, error) {
    user := models.User{
        Username: username,
        Password: password,
    }
    if err := config.DB.Create(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
