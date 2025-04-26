package accountHolderHandler

import (
	"github.com/jetaimejeteveux/simple-banking-app/internal/delivery/http/handler"
	"github.com/jetaimejeteveux/simple-banking-app/internal/service"
	"go.uber.org/zap"
)

type AccountHolderHandler struct {
	accountHolderService service.IAccountHolderService
	log                  *zap.Logger
}

func New(accountHolderService service.IAccountHolderService, log *zap.Logger) handler.IAccountHolderHandler {
	return &AccountHolderHandler{
		accountHolderService: accountHolderService,
		log:                  log,
	}
}
