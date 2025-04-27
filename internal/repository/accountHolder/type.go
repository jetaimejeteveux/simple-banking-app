package accountHolderRepository

import (
	"gorm.io/gorm"
)

type AccountHolderRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) IAccountHolderRepository {
	return &AccountHolderRepository{
		db: db,
	}
}
