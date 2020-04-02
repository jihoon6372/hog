package model

import "github.com/jinzhu/gorm"

// AuthUser 사용자 정보
type AuthUser struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);"`
	Email    string `gorm:"type:varchar(255);unique_index;not null"`
	Password string `gorm:"type:varchar(78);"`
}
