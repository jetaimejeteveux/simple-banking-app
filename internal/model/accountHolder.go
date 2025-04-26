package model

import (
	"gorm.io/gorm"
)

type AccountHolder struct {
	gorm.Model
	FullName       string  `gorm:"not null"`
	IdentityNumber string  `gorm:"uniqueIndex;not null"`
	PhoneNumber    string  `gorm:"uniqueIndex;not null"`
	AccountNumber  string  `gorm:"uniqueIndex;not null"`
	Balance        float64 `gorm:"not null;default:0"`
}

type RegisterAccountRequest struct {
	FullName       string `json:"nama" validate:"required"`
	IdentityNumber string `json:"nik" validate:"required"`
	PhoneNumber    string `json:"no_hp" validate:"required"`
}

type RegisterAccountResponse struct {
	AccountNumber string `json:"no_rekening"`
}
