package accountHolderHandler

import (
	"github.com/go-playground/validator/v10"
	"github.com/jetaimejeteveux/simple-banking-app/internal/delivery/http/handler"
	accountHolderService "github.com/jetaimejeteveux/simple-banking-app/internal/service/accountHolder"
	"go.uber.org/zap"
)

type AccountHolderHandler struct {
	accountHolderService accountHolderService.IAccountHolderService
	log                  *zap.Logger
	validator            *validator.Validate
}

func New(accountHolderService accountHolderService.IAccountHolderService, log *zap.Logger) handler.IAccountHolderHandler {
	return &AccountHolderHandler{
		accountHolderService: accountHolderService,
		log:                  log,
		validator:            validator.New(),
	}
}
