package accountHolderService

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	accountHolderRepository "github.com/jetaimejeteveux/simple-banking-app/internal/repository/accountHolder"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAccountHolderService_RegisterAccount(t *testing.T) {
	// Create controller for mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockRepo := accountHolderRepository.NewMockIAccountHolderRepository(ctrl)

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create test cases
	type args struct {
		ctx     context.Context
		request *model.RegisterAccountRequest
	}

	tests := []struct {
		name          string
		args          args
		mockSetup     func()
		wantAccountNo bool // Whether we expect a valid account number in the response
		wantErr       bool
		expectedErr   string
	}{
		{
			name: "Success - Valid Registration",
			args: args{
				ctx: context.Background(),
				request: &model.RegisterAccountRequest{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
				},
			},
			mockSetup: func() {
				// Expect check for existing phone/identity
				mockRepo.EXPECT().
					IsPhoneOrIdentityExist(gomock.Any(), "08123456789", "1234567890").
					Return(false, nil)

				// Expect account registration
				mockRepo.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantAccountNo: true,
			wantErr:       false,
		},
		{
			name: "Error - Identity or Phone Number Already Exists",
			args: args{
				ctx: context.Background(),
				request: &model.RegisterAccountRequest{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
				},
			},
			mockSetup: func() {
				// Expect check for existing phone/identity returns true (exists)
				mockRepo.EXPECT().
					IsPhoneOrIdentityExist(gomock.Any(), "08123456789", "1234567890").
					Return(true, nil)
			},
			wantAccountNo: false,
			wantErr:       true,
			expectedErr:   constants.PhoneOrIdentityExistsError,
		},
		{
			name: "Error - Checking Existing Identity/Phone Fails",
			args: args{
				ctx: context.Background(),
				request: &model.RegisterAccountRequest{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
				},
			},
			mockSetup: func() {
				// Expect check for existing phone/identity returns error
				mockRepo.EXPECT().
					IsPhoneOrIdentityExist(gomock.Any(), "08123456789", "1234567890").
					Return(false, errors.New("database error"))
			},
			wantAccountNo: false,
			wantErr:       true,
			expectedErr:   constants.RegisterAccountError,
		},
		{
			name: "Error - Registration Failed",
			args: args{
				ctx: context.Background(),
				request: &model.RegisterAccountRequest{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
				},
			},
			mockSetup: func() {
				// Expect check for existing phone/identity
				mockRepo.EXPECT().
					IsPhoneOrIdentityExist(gomock.Any(), "08123456789", "1234567890").
					Return(false, nil)

				// Expect account registration fails
				mockRepo.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(errors.New("registration failed"))
			},
			wantAccountNo: false,
			wantErr:       true,
			expectedErr:   constants.RegisterAccountError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks for this test case
			tt.mockSetup()

			// Create service with mock repository
			s := &AccountHolderService{
				accountHolderRepo: mockRepo,
				log:               logger,
			}

			// Call the function being tested
			response, err := s.RegisterAccount(tt.args.ctx, tt.args.request)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Equal(t, tt.expectedErr, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Check response
			if tt.wantAccountNo {
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.AccountNumber)

				// Check account number format (should be xxxx-xxxx-xxxx)
				parts := strings.Split(response.AccountNumber, "-")
				assert.Equal(t, 3, len(parts))
				for _, part := range parts {
					assert.Equal(t, 4, len(part))
				}
			} else if !tt.wantErr {
				assert.Fail(t, "Expected error but got response", response)
			}
		})
	}
}

// Test the account number generation function separately
func TestAccountHolderService_generateAccountNumber(t *testing.T) {
	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create service
	s := &AccountHolderService{
		log: logger,
	}

	// Generate multiple account numbers to ensure they match the format
	for i := 0; i < 10; i++ {
		accountNumber := s.generateAccountNumber()

		// Check format (xxxx-xxxx-xxxx)
		parts := strings.Split(accountNumber, "-")
		assert.Equal(t, 3, len(parts))
		for _, part := range parts {
			assert.Equal(t, 4, len(part))
		}
	}
}
