package route

import (
	"github.com/gofiber/fiber/v2"
	handler "github.com/jetaimejeteveux/simple-banking-app/internal/delivery/http/handler"
)

type RouteConfig struct {
	App                    *fiber.App
	V1AccountHolderHandler handler.IAccountHolderHandler
}

func Setup(config *RouteConfig) {
	// Register account holder routes
	config.App.Post("/v1/daftar", config.V1AccountHolderHandler.RegisterAccount)
}
