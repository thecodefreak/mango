package commands

import (
	"github.com/spf13/cobra"
	"github.com/thecodefreak/mango/internal/cli"
)

var staticPageCmd = &cobra.Command{
	Use:   "static-page <path> <file/dir>",
	Short: "Publish a static page using specified path",
	Long: `Publish a static page by providing it's path and files to publised

Examples:
  mango static-page diy-code/app1
  mango static-page pages/about.html`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		file := args[1]
		staticPage := cli.NewStaticPage(cfg, path, file)
		return staticPage.Publish()
	},
}

func init() {
	rootCmd.AddCommand(staticPageCmd)
}
