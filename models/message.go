package models

import (
    "github.com/jinzhu/gorm"
)

type Message struct {
    gorm.Model
    UserID    uint
    Content   string
    Response  string
    User      User `gorm:"foreignkey:UserID"`
}
