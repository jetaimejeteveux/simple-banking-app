package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jetaimejeteveux/simple-banking-app/internal/config"
	route "github.com/jetaimejeteveux/simple-banking-app/internal/delivery/http"
	accountHolderHandler "github.com/jetaimejeteveux/simple-banking-app/internal/delivery/http/handler/accountHolder"
	"github.com/jetaimejeteveux/simple-banking-app/internal/repository"
	accountHolderRepository "github.com/jetaimejeteveux/simple-banking-app/internal/repository/accountHolder"
	accountHolderService "github.com/jetaimejeteveux/simple-banking-app/internal/service/accountHolder"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server for the banking application.`,
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&cfg.Host, "host", "0.0.0.0", "Host to bind the server to")
	serveCmd.Flags().StringVar(&cfg.Port, "port", "8080", "Port to bind the server to")

}

func runServe(cmd *cobra.Command, args []string) {
	log.Info("Starting HTTP server...",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
	)

	// Database connection
	dsn := config.GetDatabaseDSN()
	log.Info("Connecting to database",
		zap.String("host", config.GetEnv("DB_HOST")),
		zap.String("name", config.GetEnv("DB_NAME")),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto Migrate models
	log.Info("Running database migrations")
	if err := repository.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}

	log.Info("Database connection established and migrations successfully done")

	// Init repository layer
	accountHolderRepo := accountHolderRepository.New(db)

	// Init service layer
	accountHolderSvc := accountHolderService.New(accountHolderRepo, log)

	// HTTP Handler layer
	accountHolderHandler := accountHolderHandler.New(accountHolderSvc, log)

	// Middleware
	// middlewares := middleware.NewMiddleware(log)

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Error("Unhandled error", zap.Error(err), zap.String("path", c.Path()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Middleware

	// Setup routes
	// Setup routes
	routeConfig := &route.RouteConfig{
		App:                    app,
		V1AccountHolderHandler: accountHolderHandler,
	}
	route.Setup(routeConfig)
	log.Info("Routes initialized")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info("Shutting down server...")

		// Add a timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
		defer cancel()

		// Get the database connection instance
		sqlDB, err := db.DB()
		if err != nil {
			log.Error("Failed to get database instance", zap.Error(err))
		} else {
			log.Info("Closing database connection")
			sqlDB.Close()
		}

		// Shutdown the server
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatal("Server shutdown failed", zap.Error(err))
		}

		log.Info("Server shutdown complete")
		os.Exit(0)
	}()

	// Start server
	serverAddress := cfg.Host + ":" + cfg.Port
	log.Info("Server is ready to handle requests", zap.String("address", serverAddress))
	if err := app.Listen(serverAddress); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
