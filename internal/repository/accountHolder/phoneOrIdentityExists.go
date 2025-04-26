package accountHolderRepository

import (
	"context"

	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
)

func (r *AccountHolderRepository) IsPhoneOrIdentityExist(ctx context.Context, phoneNumber, identityNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.AccountHolder{}).
		Where("identity_number = ? OR phone_number = ?", phoneNumber, identityNumber).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
