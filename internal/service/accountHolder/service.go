package accountHolderService

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=accountHolderService
type IAccountHolderService interface {
	RegisterAccount(ctx context.Context, request *model.RegisterAccountRequest) (*model.RegisterAccountResponse, error)
	Deposit(ctx context.Context, request *model.DepositRequest) (*model.DepositResponse, error)
	Withdraw(ctx context.Context, request *model.WithdrawRequest) (*model.WithdrawResponse, error)
	GetBalance(ctx context.Context, request *model.GetBalanceRequest) (*model.GetBalanceResponse, error)
}
