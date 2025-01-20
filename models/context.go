package models

import "github.com/jinzhu/gorm"

type Context struct {
    gorm.Model
    UserID uint   `json:"user_id"`
    Key    string `json:"key"`   
    Value  string `json:"value"`
}
