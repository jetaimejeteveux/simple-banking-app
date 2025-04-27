package accountHolderService

import (
	"context"
	"errors"
	"testing"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	accountHolderRepository "github.com/jetaimejeteveux/simple-banking-app/internal/repository/accountHolder"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestAccountHolderService_Deposit(t *testing.T) {
	// Create controller for mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockRepo := accountHolderRepository.NewMockIAccountHolderRepository(ctrl)

	// Create logger
	logger, _ := zap.NewDevelopment()

	type args struct {
		ctx     context.Context
		request *model.DepositRequest
	}

	tests := []struct {
		name        string
		args        args
		mockSetup   func()
		wantBalance float64
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Success - Valid Deposit",
			args: args{
				ctx: context.Background(),
				request: &model.DepositRequest{
					AccountNumber: "1234-5678-9012",
					Amount:        5000,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetByAccountNumber(gomock.Any(), "1234-5678-9012").
					Return(&model.AccountHolder{
						AccountNumber: "1234-5678-9012",
						Balance:       10000,
					}, nil).Times(1)

				mockRepo.EXPECT().
					UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			wantBalance: 15000,
			wantErr:     false,
		},
		{
			name: "Error - Account Not Found",
			args: args{
				ctx: context.Background(),
				request: &model.DepositRequest{
					AccountNumber: "1234-5678-9012",
					Amount:        5000,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetByAccountNumber(gomock.Any(), "1234-5678-9012").
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantBalance: 0,
			wantErr:     true,
			expectedErr: constants.AccountNotFoundError,
		},
		{
			name: "Error - Fetching Account Fails",
			args: args{
				ctx: context.Background(),
				request: &model.DepositRequest{
					AccountNumber: "1234-5678-9012",
					Amount:        5000,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetByAccountNumber(gomock.Any(), "1234-5678-9012").
					Return(nil, errors.New("database error"))
			},
			wantBalance: 0,
			wantErr:     true,
			expectedErr: constants.FetchAccountHolderError,
		},
		{
			name: "Error - Updating Balance Fails",
			args: args{
				ctx: context.Background(),
				request: &model.DepositRequest{
					AccountNumber: "1234-5678-9012",
					Amount:        5000,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetByAccountNumber(gomock.Any(), "1234-5678-9012").
					Return(&model.AccountHolder{
						AccountNumber: "1234-5678-9012",
						Balance:       10000,
					}, nil)

				mockRepo.EXPECT().
					UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("update failed"))
			},
			wantBalance: 0,
			wantErr:     true,
			expectedErr: constants.DepositError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := &AccountHolderService{
				accountHolderRepo: mockRepo,
				log:               logger,
			}

			resp, err := s.Deposit(tt.args.ctx, tt.args.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Equal(t, tt.expectedErr, err.Error())
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantBalance, resp.Balance)
			}
		})
	}
}
