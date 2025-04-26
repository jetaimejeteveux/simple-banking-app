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
	FullName       string `json:"full_name"`
	IdentityNumber string `json:"identity_number"`
	PhoneNumber    string `json:"phone_number"`
}

type RegisterAccountResponse struct {
	AccountNumber string `json:"account_number"`
}
