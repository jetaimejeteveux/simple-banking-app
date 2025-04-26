package accountHolderHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"go.uber.org/zap"
)

func (h *AccountHolderHandler) RegisterAccount(c *fiber.Ctx) error {
	var req model.RegisterAccountRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": constants.InvalidRequestError,
		})
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": constants.MissingFieldError,
		})
	}

	response, err := h.accountHolderService.RegisterAccount(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
