package handler

import (
	"github.com/gofiber/fiber/v2"
)

type IAccountHolderHandler interface {
	RegisterAccount(ctx *fiber.Ctx) error
}
