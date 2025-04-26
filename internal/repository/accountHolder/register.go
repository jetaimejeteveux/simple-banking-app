package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) Register(ctx context.Context, accountHolder *model.AccountHolder) error {
	if err := r.db.WithContext(ctx).Create(accountHolder).Error; err != nil {
		return err
	}
	return nil
}
