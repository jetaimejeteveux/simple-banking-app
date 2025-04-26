package accountHolderRepository

import (
	"github.com/jetaimejeteveux/simple-banking-app/internal/repository"
	"gorm.io/gorm"
)

type AccountHolderRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.IAccountHolderRepository {
	return &AccountHolderRepository{
		db: db,
	}
}
