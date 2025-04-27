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

func TestAccountHolderHandler_RegisterAccount(t *testing.T) {
	// Setup
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
			name: "Success - Valid Registration",
			requestBody: model.RegisterAccountRequest{
				FullName:       "John Doe",
				IdentityNumber: "1234567890",
				PhoneNumber:    "08123456789",
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
					RegisterAccount(gomock.Any(), gomock.Any()).
					Return(&model.RegisterAccountResponse{
						AccountNumber: "1234-5678-9012",
					}, nil)
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"no_rekening": "1234-5678-9012",
			},
		},
		{
			name: "Error - Invalid Request Body",
			requestBody: `{
				"invalid json"
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
			requestBody: model.RegisterAccountRequest{
				// Missing all required fields
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
			requestBody: model.RegisterAccountRequest{
				FullName:       "John Doe",
				IdentityNumber: "1234567890",
				PhoneNumber:    "08123456789",
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
					RegisterAccount(gomock.Any(), gomock.Any()).
					Return(nil, errors.New(constants.PhoneOrIdentityExistsError))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"remark": constants.PhoneOrIdentityExistsError,
			},
		},
		{
			name: "Error - Partial Request Body",
			requestBody: model.RegisterAccountRequest{
				FullName: "John Doe",
				// Missing IdentityNumber and PhoneNumber
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service expectations
			tt.mockServiceSetup()

			// Create handler with mock service
			handler := tt.handler()

			// Create a new Fiber app
			app := fiber.New(fiber.Config{
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					// Custom error handler to avoid test panic on internal errors
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"remark": "internal server error",
					})
				},
			})
			app.Post("/v1/daftar", handler.RegisterAccount)

			// Create request body
			var reqBody []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/v1/daftar", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			// Check response body
			if tt.expectedBody != nil {
				// Read response body
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				// Parse JSON
				var responseMap map[string]interface{}
				err = json.Unmarshal(body, &responseMap)
				assert.NoError(t, err)

				// Check expected fields
				for key, expectedValue := range tt.expectedBody {
					assert.Equal(t, expectedValue, responseMap[key])
				}
			}
		})
	}
}
