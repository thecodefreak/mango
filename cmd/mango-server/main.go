package main

import (
	"fmt"
	"net"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadServerConf()
		if err != nil {
			return fmt.Errorf("Config load failed: %w", err)
		}

		if ServerConfig.Addr != "" {
			config.Addr = ServerConfig.Addr
		}

		if ServerConfig.DocumentRoot != "" {
			config.DocumentRoot = ServerConfig.DocumentRoot
		}

		err = validateAddrAndDocumentRoot(config)
		if err != nil {
			return fmt.Errorf("Invalid Config: %w", err)
		}

		err = server.InitServer(config)
		return err
	},
}

func validateAddrAndDocumentRoot(c *config.ServerConfig) error {
	if c == nil {
		return fmt.Errorf("Config not loaded.")
	}

	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("Server address is required")
	}

	if _, _, err := net.SplitHostPort(c.Addr); err != nil {
		return fmt.Errorf("invalid server address %q, must be in host:port format: %w", c.Addr, err)
	}

	if strings.TrimSpace(c.DocumentRoot) == "" {
		return fmt.Errorf("document root is required")
	}

	if err := helpers.IsDirWritable(c.DocumentRoot); err != nil {
		return fmt.Errorf("document root is not writable %q: %w", c.DocumentRoot, err)
	}

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates config file for mango-server",
	RunE: func(c *cobra.Command, args []string) error {
		v := config.InitConfig()
		err := config.CreateConfig(v, "server")
		if err != nil {
			return fmt.Errorf("Config init failed: %w", err)
		}
		return nil
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&ServerConfig.Addr, "server-addr", "", "Address to bind to")
	rootCmd.Flags().StringVar(&ServerConfig.DocumentRoot, "document-root", "", "Document root for static files")
}
