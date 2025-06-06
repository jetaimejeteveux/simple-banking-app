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

type DepositRequest struct {
	AccountNumber string  `json:"no_rekening" validate:"required"`
	Amount        float64 `json:"nominal" validate:"required"`
}
type DepositResponse struct {
	Balance float64 `json:"saldo"`
}

type WithdrawRequest struct {
	AccountNumber string  `json:"no_rekening" validate:"required"`
	Amount        float64 `json:"nominal" validate:"required"`
}

type WithdrawResponse struct {
	Balance float64 `json:"saldo"`
}

type GetBalanceRequest struct {
	AccountNumber string `validate:"required"`
}

type GetBalanceResponse struct {
	Balance float64 `json:"saldo"`
}
