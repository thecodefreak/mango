package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thecodefreak/mango/internal/config"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "mango",
	Short: "Multi swiss army knife for developers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Welcome to Mango 🥭! \n\nUse --help to see available commands.\n")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("Unable to load config, %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() error {
	return rootCmd.Execute()
}
