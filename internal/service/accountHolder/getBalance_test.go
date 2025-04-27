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

func TestAccountHolderService_GetBalance(t *testing.T) {
	// Create controller for mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockRepo := accountHolderRepository.NewMockIAccountHolderRepository(ctrl)

	// Create logger
	logger, _ := zap.NewDevelopment()

	type args struct {
		ctx     context.Context
		request *model.GetBalanceRequest
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
			name: "Success - Get Balance",
			args: args{
				ctx: context.Background(),
				request: &model.GetBalanceRequest{
					AccountNumber: "123456789012",
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetByAccountNumber(gomock.Any(), "1234-5678-9012").
					Return(&model.AccountHolder{
						AccountNumber: "1234-5678-9012",
						Balance:       10000,
					}, nil).Times(1)
			},
			wantBalance: 10000,
			wantErr:     false,
		},
		{
			name: "Error - Account Not Found",
			args: args{
				ctx: context.Background(),
				request: &model.GetBalanceRequest{
					AccountNumber: "123456789012",
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
				request: &model.GetBalanceRequest{
					AccountNumber: "123456789012",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			s := &AccountHolderService{
				accountHolderRepo: mockRepo,
				log:               logger,
			}

			resp, err := s.GetBalance(tt.args.ctx, tt.args.request)

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
