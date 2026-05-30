package config

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if (line == "" || strings.HasPrefix(line, "#")) {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key != "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func Get(key string) string {
	return os.Getenv(key)
}

func MustGet(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}