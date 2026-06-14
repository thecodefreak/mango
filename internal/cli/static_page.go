package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thecodefreak/mango/internal/api"
	"github.com/thecodefreak/mango/internal/config"
	"github.com/thecodefreak/mango/internal/handlers"
	"github.com/thecodefreak/mango/internal/helpers"
)

type StaticPage struct {
	Cfg    *config.Config
	Client *api.Client
	Path   string
	File   string
}

func NewStaticPage(cfg *config.Config, path, file string) *StaticPage {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		return nil
	}

	return &StaticPage{
		Cfg:    cfg,
		Client: api.NewClient(cfg.Server, cfg.ApiToken),
		Path:   strings.TrimSpace(path),
		File:   strings.TrimSpace(file),
	}
}

func (s *StaticPage) Publish() error {
	if s.Path == "" || s.File == "" {
		return fmt.Errorf("path and file must be provided")
	}

	client := api.NewClient(s.Cfg.Server, s.Cfg.ApiToken)

	fileInfo, err := os.Stat(s.File)

	absPath, _ := filepath.Abs(s.File)
	if err != nil {
		return fmt.Errorf("Given file or directory does not exist: %s\n", s.File)
	}

	if fileInfo.IsDir() {
		files := []handlers.StaticPageContentFiles{}
		err := filepath.WalkDir(s.File, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("Error occurred while walking directory %s: %s", s.File, err)
			}
			if d.IsDir() {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("Failed to read file %s", path)
			}

			fileHash, err := helpers.FileChecksum(path)
			if err != nil {
				return fmt.Errorf("Failed to calculate checksum for file %s", path)
			}
			files = append(files, handlers.StaticPageContentFiles{
				Path:    strings.TrimPrefix(path, absPath),
				Content: base64.StdEncoding.EncodeToString(content),
				Hash:    fileHash,
			})

			return nil
		})

		if err != nil {
			return fmt.Errorf("Failed while walking through directory %s: %s", s.File, err)
		}

		err = client.StaticPageCreate(handlers.StaticPageContent{
			PagePath: s.Path,
			Files:    files,
		})

		if err != nil {
			return fmt.Errorf("Failed to create or update static page. \n\n%w\n", err)
		}

		fmt.Printf("Successfully published static page at path: %s\n", s.Path)
	} else {
		fmt.Printf("File info for %s:\n", s.File)
		fmt.Printf("  Type: File\n")
	}

	return nil
}

func UploadFiles() {

}
