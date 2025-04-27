package accountHolderHandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	accountHolderService "github.com/jetaimejeteveux/simple-banking-app/internal/service/accountHolder"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/constants"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAccountHolderHandler_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := accountHolderService.NewMockIAccountHolderService(ctrl)
	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	tests := []struct {
		name               string
		requestBody        interface{}
		handler            func() *AccountHolderHandler
		mockServiceSetup   func()
		expectedStatusCode int
		expectedBody       map[string]interface{}
	}{
		{
			name: "Success - Valid Deposit",
			requestBody: model.DepositRequest{
				AccountNumber: "1234-5678-9012",
				Amount:        100000,
			},
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
					validator:            validate,
				}
			},
			mockServiceSetup: func() {
				mockService.EXPECT().
					Deposit(gomock.Any(), gomock.Any()).
					Return(&model.DepositResponse{
						Balance: 200000,
					}, nil)
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"saldo": float64(200000),
			},
		},
		{
			name: "Error - Invalid Request Body",
			requestBody: `{
				invalid json
			}`,
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
					validator:            validate,
				}
			},
			mockServiceSetup:   func() {},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"remark": constants.InvalidRequestError,
			},
		},
		{
			name:        "Error - Missing Required Fields",
			requestBody: model.DepositRequest{
				// Missing fields: AccountNumber, Amount
			},
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
					validator:            validate,
				}
			},
			mockServiceSetup:   func() {},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"remark": constants.MissingFieldError,
			},
		},
		{
			name: "Error - Service Error",
			requestBody: model.DepositRequest{
				AccountNumber: "1234-5678-9012",
				Amount:        50000,
			},
			handler: func() *AccountHolderHandler {
				return &AccountHolderHandler{
					accountHolderService: mockService,
					log:                  logger,
					validator:            validate,
				}
			},
			mockServiceSetup: func() {
				mockService.EXPECT().
					Deposit(gomock.Any(), gomock.Any()).
					Return(nil, errors.New(constants.AccountNotFoundError))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"remark": constants.AccountNotFoundError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockServiceSetup()

			handler := tt.handler()

			app := fiber.New(fiber.Config{
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"remark": "internal server error",
					})
				},
			})
			app.Post("/v1/tabung", handler.Deposit)

			var reqBody []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/v1/tabung", bytes.NewReader(reqBody))
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
