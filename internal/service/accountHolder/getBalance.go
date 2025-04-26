package accountHolderService

import (
	"context"
	"errors"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
	"go.uber.org/zap"
)

func (s *AccountHolderService) GetBalance(ctx context.Context, request *model.GetBalanceRequest) (*model.GetBalanceResponse, error) {
	logger := s.log.With(
		zap.String("Service", "accountHolder"),
		zap.String("method", "GetBalance"),
	)

	accountHolder, err := s.accountHolderRepo.GetByAccountNumber(ctx, request.AccountNumber)
	if err != nil {
		if helper.IsRecordNotFound(err) {
			logger.Error("Account not found", zap.String("AccountNumber", request.AccountNumber))
			return nil, errors.New("account not found")
		}
		logger.Error("Error fetching account holder", zap.Error(err))
		return nil, err
	}

	return &model.GetBalanceResponse{
		Balance: accountHolder.Balance,
	}, nil
}
