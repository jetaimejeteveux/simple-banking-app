package repository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

type IAccountHolderRepository interface {
	Register(ctx context.Context, accountHolder *model.AccountHolder) error
	GetByIdentityNumber(ctx context.Context, identityNumber string) (*model.AccountHolder, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.AccountHolder, error)
}
