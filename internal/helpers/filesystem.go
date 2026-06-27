package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func FileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirWritable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access directory %q: %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", path)
	}

	f, err := os.CreateTemp(path, ".write-test-*")
	if err != nil {
		return fmt.Errorf("directory %q is not writable: %w", path, err)
	}

	name := f.Name()

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close temp file %q: %w", name, err)
	}

	if err := os.Remove(name); err != nil {
		return fmt.Errorf("failed to remove temp file %q: %w", name, err)
	}

	return nil
}

func PathWoSlash(path string, isPrefix bool) string {
	if len(path) > 0 {
		if isPrefix && !strings.HasSuffix(path, "/") {
			path += "/"
		} else if !isPrefix && strings.HasSuffix(path, "/") {
			path = strings.TrimSuffix(path, "/")
		}
	}
	return path
}
