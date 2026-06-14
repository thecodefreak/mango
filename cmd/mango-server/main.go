package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thecodefreak/mango/internal/config"
	"github.com/thecodefreak/mango/internal/helpers"
	"github.com/thecodefreak/mango/internal/server"
)

var ServerConfig = config.ServerConfig{}

var rootCmd = &cobra.Command{
	Use:   "mango-server",
	Short: "Server for mango cli",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadServerConf()
		if err != nil {
			fmt.Print(err)
		}

		if ServerConfig.Addr != "" {
			config.Addr = ServerConfig.Addr
		}

		if ServerConfig.DocumentRoot != "" {
			config.DocumentRoot = ServerConfig.DocumentRoot
		}

		err = validateAddrAndDocumentRoot(config)
		if err != nil {
			fmt.Printf("%s \n", err)
			os.Exit(1)
		}

		server.InitServer(config)
	},
}

func validateAddrAndDocumentRoot(c *config.ServerConfig) error {
	if c.Addr != "" && !strings.Contains(c.Addr, ":") {
		return fmt.Errorf("Invalid server address: %s. Must be in the format host:port", c.Addr)
	}
	if c.DocumentRoot != "" {
		err := helpers.IsDirWritable(c.DocumentRoot)
		if err != nil {
			return fmt.Errorf("Document root is not writable: %s", c.DocumentRoot)
		}
	}
	return nil
}

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&ServerConfig.Addr, "server-addr", "", "Address to bind to")
	rootCmd.Flags().StringVar(&ServerConfig.DocumentRoot, "document-root", "", "Document root for static files")
}
