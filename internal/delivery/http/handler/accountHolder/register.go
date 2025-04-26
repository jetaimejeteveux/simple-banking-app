package accountHolderHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"go.uber.org/zap"
)

func (h *AccountHolderHandler) RegisterAccount(c *fiber.Ctx) error {
	var req model.RegisterAccountRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	response, err := h.accountHolderService.RegisterAccount(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account_number": response.AccountNumber,
	})
}
