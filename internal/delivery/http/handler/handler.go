package handler

import (
	"github.com/gofiber/fiber/v2"
)

type IAccountHolderHandler interface {
	RegisterAccount(ctx *fiber.Ctx) error
	Deposit(ctx *fiber.Ctx) error
	Withdraw(ctx *fiber.Ctx) error
	GetBalance(ctx *fiber.Ctx) error
}
