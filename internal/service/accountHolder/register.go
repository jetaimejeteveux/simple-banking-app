package accountHolderService

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"go.uber.org/zap"
)

func (s *AccountHolderService) RegisterAccount(ctx context.Context, request *model.RegisterAccountRequest) (*model.RegisterAccountResponse, error) {
	logger := s.log.With(
		zap.String("Service", "accountHolder"),
		zap.String("method", "RegisterAccount"),
	)

	accountNumber := s.generateAccountNumber()
	accountHolder := &model.AccountHolder{
		FullName:       request.FullName,
		IdentityNumber: request.IdentityNumber,
		PhoneNumber:    request.PhoneNumber,
		AccountNumber:  accountNumber,
		Balance:        0,
	}

	// Check if IdentityNumber already exists
	exists, err := s.accountHolderRepo.IsPhoneOrIdentityExist(ctx, request.PhoneNumber, request.IdentityNumber)
	if err != nil {
		logger.Error("Error checking existing account", zap.Error(err))
		return nil, errors.New(constants.RegisterAccountError)
	}
	if exists {
		logger.Warn("Attempt to register with existing IdentityNumber or PhoneNumber",
			zap.String("IdentityNumber", request.IdentityNumber),
			zap.String("PhoneNumber", request.PhoneNumber),
		)
		return nil, errors.New(constants.PhoneOrIdentityExistsError)
	}

	err = s.accountHolderRepo.Register(ctx, accountHolder)
	if err != nil {
		logger.Error("Error registering account holder", zap.Error(err))
		return nil, errors.New(constants.RegisterAccountError)
	}

	return &model.RegisterAccountResponse{
		AccountNumber: accountNumber,
	}, nil

}

func (s *AccountHolderService) generateAccountNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a 12 digit random number with 4 - 4 - 4 format
	return fmt.Sprintf("%04d-%04d-%04d", r.Intn(10000), r.Intn(10000), r.Intn(10000))
}
