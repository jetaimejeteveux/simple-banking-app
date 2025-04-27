package accountHolderService

import (
	"context"
	"errors"
	"fmt"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
	"go.uber.org/zap"
)

func (s *AccountHolderService) GetBalance(ctx context.Context, request *model.GetBalanceRequest) (*model.GetBalanceResponse, error) {
	logger := s.log.With(
		zap.String("Service", "accountHolder"),
		zap.String("method", "GetBalance"),
	)

	formattedAccountNumber := fmt.Sprintf("%s-%s-%s",
		request.AccountNumber[0:4],
		request.AccountNumber[4:8],
		request.AccountNumber[8:12])

	accountHolder, err := s.accountHolderRepo.GetByAccountNumber(ctx, formattedAccountNumber)
	if err != nil {
		if helper.IsRecordNotFound(err) {
			logger.Error("Account not found", zap.String("AccountNumber", request.AccountNumber))
			return nil, errors.New(constants.AccountNotFoundError)
		}
		logger.Error("Error fetching account holder", zap.Error(err))
		return nil, errors.New(constants.FetchAccountHolderError)
	}

	return &model.GetBalanceResponse{
		Balance: accountHolder.Balance,
	}, nil
}
