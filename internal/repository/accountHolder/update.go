package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) Update(ctx context.Context, accountHolder *model.AccountHolder) error {
	return r.db.WithContext(ctx).Save(accountHolder).Error
}
