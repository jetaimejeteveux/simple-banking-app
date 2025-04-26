package accountHolderService

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/helper"
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
	existingAccountHolder, err := s.accountHolderRepo.GetByIdentityNumber(ctx, request.IdentityNumber)
	if err != nil && !helper.IsRecordNotFound(err) {
		logger.Error("Error checking IdentityNumber", zap.Error(err))
		return nil, err
	}
	if existingAccountHolder != nil {
		logger.Error("IdentityNumber already exists", zap.String("IdentityNumber", request.IdentityNumber))
		return nil, errors.New("IdentityNumber already exists")
	}

	// Check if PhoneNumber already exists
	existingAccountHolder, err = s.accountHolderRepo.GetByPhoneNumber(ctx, request.PhoneNumber)
	if err != nil && !helper.IsRecordNotFound(err) {
		logger.Error("Error checking PhoneNumber", zap.Error(err))
		return nil, err
	}
	if existingAccountHolder != nil {
		logger.Error("PhoneNumber already exists", zap.String("PhoneNumber", request.PhoneNumber))
		return nil, errors.New("PhoneNumber already exists")

	}

	err = s.accountHolderRepo.Register(ctx, accountHolder)
	if err != nil {
		logger.Error("Error registering account holder", zap.Error(err))
		return nil, err
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
