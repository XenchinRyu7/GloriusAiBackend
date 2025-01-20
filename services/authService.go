package services

import (
	"gloriusaiapi/models"
	"gloriusaiapi/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret")

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func AuthenticateUser(username, password string) (*models.User, string, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", err
	}

	token, err := GenerateToken(user.ID)
	return user, token, err
}

func RegisterUser(db *gorm.DB, name, email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: name,
		Password: string(hashedPassword),
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
