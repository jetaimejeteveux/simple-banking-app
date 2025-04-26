package service

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

type IAccountHolderService interface {
	RegisterAccount(ctx context.Context, request *model.RegisterAccountRequest) (*model.RegisterAccountResponse, error)
	Deposit(ctx context.Context, request *model.DepositRequest) (*model.DepositResponse, error)
}
