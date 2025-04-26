package accountHolderHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
)

func (h *AccountHolderHandler) GetBalance(ctx *fiber.Ctx) error {
	accountNumber := ctx.Params("no_rekening")
	if accountNumber == "" {
		h.log.Warn("Account number is missing in path")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Account number is required in the path",
		})
	}

	request := model.GetBalanceRequest{
		AccountNumber: accountNumber,
	}

	response, err := h.accountHolderService.GetBalance(ctx.Context(), &request)
	if err != nil {
		if helper.IsRecordNotFound(err) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"remark": "Account not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Failed to get account",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
