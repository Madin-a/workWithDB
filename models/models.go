package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `json:"-"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required" validate:"min=1,max=16" `
}

type LogPas struct {
	gorm.Model
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
