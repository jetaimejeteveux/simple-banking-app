package accountHolderHandler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	accountHolderService "github.com/jetaimejeteveux/simple-banking-app/internal/service/accountHolder"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAccountHolderHandler_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := accountHolderService.NewMockIAccountHolderService(ctrl)
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name               string
		accountNumber      string
		handler            func() *AccountHolderHandler
		mockServiceSetup   func()
		expectedStatusCode int
		expectedBody       map[string]interface{}
	}{
		{
			name:          "Success - Valid Account Number",
			accountNumber: "123456789012",
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
				}
			},
			mockServiceSetup: func() {
				mockService.EXPECT().
					GetBalance(gomock.Any(), &model.GetBalanceRequest{
						AccountNumber: "123456789012",
					}).
					Return(&model.GetBalanceResponse{
						Balance: 500000,
					}, nil)
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"saldo": float64(500000),
			},
		},
		{
			name:          "Error - Missing Account Number",
			accountNumber: "",
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
				}
			},
			mockServiceSetup:   func() {},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"remark": constants.MissingFieldError,
			},
		},
		{
			name:          "Error - Account Not Found",
			accountNumber: "999999999999",
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
				}
			},
			mockServiceSetup: func() {
				mockService.EXPECT().
					GetBalance(gomock.Any(), gomock.Any()).
					Return(nil, errors.New(constants.AccountNotFoundError)) // Fixed to return an error object
			},
			expectedStatusCode: fiber.StatusNotFound,
			expectedBody: map[string]interface{}{
				"remark": constants.AccountNotFoundError,
			},
		},
		{
			name:          "Error - Service Error",
			accountNumber: "123456789012",
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
				}
			},
			mockServiceSetup: func() {
				mockService.EXPECT().
					GetBalance(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"remark": constants.FetchAccountHolderError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockServiceSetup()

			handler := tt.handler()

			app := fiber.New()
			app.Get("/v1/saldo/:no_rekening", handler.GetBalance)

			req := httptest.NewRequest("GET", "/v1/saldo/"+tt.accountNumber, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedBody != nil {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var responseMap map[string]interface{}
				err = json.Unmarshal(body, &responseMap)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					assert.Equal(t, expectedValue, responseMap[key])
				}
			}
		})
	}
}
