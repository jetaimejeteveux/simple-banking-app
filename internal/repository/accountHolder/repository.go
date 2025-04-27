package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=accountHolderRepository
type IAccountHolderRepository interface {
	Register(ctx context.Context, accountHolder *model.AccountHolder) error
	GetByIdentityNumber(ctx context.Context, identityNumber string) (*model.AccountHolder, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.AccountHolder, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*model.AccountHolder, error)
	Update(ctx context.Context, accountHolder *model.AccountHolder) error
	IsPhoneOrIdentityExist(ctx context.Context, phoneNumber, identityNumber string) (bool, error)
}
