package models

import "gorm.io/gorm"

type Usermodel struct {
	gorm.Model
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
}
