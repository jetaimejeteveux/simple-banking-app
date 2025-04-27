package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) UpdateBalance(ctx context.Context, accountNumber string, balance float64) error {
	return r.db.WithContext(ctx).
		Model(&model.AccountHolder{}).
		Where("account_number = ?", accountNumber).
		Update("balance", balance).
		Error
}
