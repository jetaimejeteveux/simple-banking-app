package accountHolderService

import (
	"context"
	"errors"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
	"go.uber.org/zap"
)

func (s *AccountHolderService) Deposit(ctx context.Context, request *model.DepositRequest) (*model.DepositResponse, error) {
	logger := s.log.With(
		zap.String("Service", "accountHolder"),
		zap.String("method", "Deposit"),
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

	accountHolder.Balance += request.Amount

	err = s.accountHolderRepo.UpdateBalance(ctx, accountHolder.AccountNumber, accountHolder.Balance)
	if err != nil {
		logger.Error("Error updating account holder balance", zap.Error(err))
		return nil, errors.New(constants.DepositError)
	}

	return &model.DepositResponse{
		Balance: accountHolder.Balance,
	}, nil
}
