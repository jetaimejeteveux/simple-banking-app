package cmd

import (
	"fmt"
	"os"

	"github.com/jetaimejeteveux/simple-banking-app/internal/config"
	"github.com/jetaimejeteveux/simple-banking-app/internal/utils/logger"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	log     = logger.NewLogger()
)

var rootCmd = &cobra.Command{
	Use:   "cobra-cli",
	Short: "A generator for Simple Banking App",
	Long:  `Simple banking app to display the ability of a banking system with basic operations like deposit, withdraw, and balance inquiry.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	var err error
	cfg, err = config.NewConfig()
	if err != nil {
		fmt.Println("Error initializing config: ", err)
		os.Exit(1)
	}
}

func initConfig() {
	if err := config.LoadEnv(); err != nil {
		fmt.Println("Warning: error loading environment variables:", err)
	}
}
