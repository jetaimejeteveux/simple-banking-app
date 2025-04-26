package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) GetByPhoneNumber(ctx context.Context, phone string) (*model.AccountHolder, error) {
	var accountHolder model.AccountHolder
	if err := r.db.WithContext(ctx).Where("phone_number = ?", phone).First(&accountHolder).Error; err != nil {
		return nil, err
	}
	return &accountHolder, nil
}
