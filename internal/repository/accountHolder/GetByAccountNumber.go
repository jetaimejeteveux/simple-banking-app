package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*model.AccountHolder, error) {
	accountHolder := &model.AccountHolder{}
	err := r.db.WithContext(ctx).Where("account_number = ?", accountNumber).First(accountHolder).Error
	if err != nil {
		return nil, err
	}
	return accountHolder, nil
}
