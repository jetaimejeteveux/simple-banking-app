package repository

import (
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.AccountHolder{},
	)
}
