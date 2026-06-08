package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/thecodefreak/mango/internal/api"
)

var client = api.NewClient()

type StaticPage struct {
	Path string
	File string
}

func NewStaticPage(path, file string) *StaticPage {
	return &StaticPage{
		Path: strings.TrimSpace(path),
		File: strings.TrimSpace(file),
	}
}

func (s *StaticPage) Publish() error {
	if s.Path == "" || s.File == "" {
		return fmt.Errorf("path and file must be provided")
	}

	fileInfo, err := os.Stat(s.File)
	if err != nil {
		fmt.Errorf("Given file or directory does not exist: %s\n", s.File)
		return err
	}

	if fileInfo.IsDir() {
		listFiles, err := os.ReadDir(s.File)
		if err != nil {
			fmt.Printf("Failed to read files from %s", s.File)
		}
		for _, file := range listFiles {
			fmt.Print(file)
		}
	} else {
		fmt.Printf("File info for %s:\n", s.File)
		fmt.Printf("  Type: File\n")
	}

	return nil
}

func UploadFiles() {

}

