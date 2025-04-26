package accountHolderService

import (
	"context"
	"errors"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
	"go.uber.org/zap"
)

func (s *AccountHolderService) Withdraw(ctx context.Context, request *model.WithdrawRequest) (*model.WithdrawResponse, error) {
	logger := s.log.With(
		zap.String("Service", "accountHolder"),
		zap.String("method", "Withdraw"),
	)

	accountHolder, err := s.accountHolderRepo.GetByAccountNumber(ctx, request.AccountNumber)
	if err != nil {
		if helper.IsRecordNotFound(err) {
			logger.Error("Account not found", zap.String("AccountNumber", request.AccountNumber))
			return nil, errors.New(constants.AccountNotFoundError)
		}
		logger.Error("Error fetching account holder", zap.Error(err))
		return nil, errors.New(constants.FetchAccountHolderError)
	}

	if accountHolder.Balance < request.Amount {
		logger.Error("Insufficient funds", zap.Float64("Balance", accountHolder.Balance), zap.Float64("WithdrawAmount", request.Amount))
		return nil, errors.New(constants.InsufficientBalanceError)
	}

	accountHolder.Balance -= request.Amount

	err = s.accountHolderRepo.Update(ctx, accountHolder)
	if err != nil {
		logger.Error("Error updating account holder balance", zap.Error(err))
		return nil, errors.New(constants.WithdrawError)
	}

	return &model.WithdrawResponse{
		Balance: accountHolder.Balance,
	}, nil
}
