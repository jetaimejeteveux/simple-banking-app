package accountHolderService

import (
	"github.com/jetaimejeteveux/simple-banking-app/internal/repository"
	"go.uber.org/zap"
)

type AccountHolderService struct {
	accountHolderRepo repository.IAccountHolderRepository
	log               *zap.Logger
}

func New(accountHolderRepo repository.IAccountHolderRepository, log *zap.Logger) *AccountHolderService {
	return &AccountHolderService{
		accountHolderRepo: accountHolderRepo,
		log:               log,
	}
}
