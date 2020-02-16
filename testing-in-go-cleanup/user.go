package main

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"type:varchar(100);column:username"`
	Password string `gorm:"type:varchar(100);column:password"`
}
