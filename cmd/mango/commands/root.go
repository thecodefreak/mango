package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "mango",
	Short: "Multi swiss army knife for developers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Mango! Use --help to see available commands.")
	},
}

var verbose bool
   
func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() error {
	return rootCmd.Execute()
}
