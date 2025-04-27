package accountHolderService

import (
	accountHolderRepository "github.com/jetaimejeteveux/simple-banking-app/internal/repository/accountHolder"
	"go.uber.org/zap"
)

type AccountHolderService struct {
	accountHolderRepo accountHolderRepository.IAccountHolderRepository
	log               *zap.Logger
}

func New(accountHolderRepo accountHolderRepository.IAccountHolderRepository, log *zap.Logger) *AccountHolderService {
	return &AccountHolderService{
		accountHolderRepo: accountHolderRepo,
		log:               log,
	}
}
