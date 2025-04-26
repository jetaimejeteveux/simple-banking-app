package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) GetByIdentityNumber(ctx context.Context, identityNumber string) (*model.AccountHolder, error) {
	var accountHolder model.AccountHolder
	if err := r.db.WithContext(ctx).Where("identity_number = ?", identityNumber).First(&accountHolder).Error; err != nil {
		return nil, err
	}
	return &accountHolder, nil
}
