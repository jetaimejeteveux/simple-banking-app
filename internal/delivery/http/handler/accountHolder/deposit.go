package accountHolderHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"go.uber.org/zap"
)

func (h *AccountHolderHandler) Deposit(ctx *fiber.Ctx) error {
	var request model.DepositRequest
	if err := ctx.BodyParser(&request); err != nil {
		h.log.Warn("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Invalid request payload",
		})
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Required fields are missing or invalid",
		})
	}

	response, err := h.accountHolderService.Deposit(ctx.Context(), &request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Failed to deposit",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
