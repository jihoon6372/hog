package model

import "github.com/jinzhu/gorm"

// User 사용자 정보
type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(100);"`
	Email     string `gorm:"type:varchar(255);unique_index;not null"`
	Password  string `gorm:"type:varchar(78);"`
	Profile   Profile
	ProfileID int
}

// Profile 사용자 프로필
type Profile struct {
	gorm.Model
	Address string
}
