package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	tokenFileName = "token.yaml"
)

func TokenFile() string {
	var once sync.Once

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot fetch required directory: %w", err))
	}

	configDir := filepath.Join(userConfigDir, "chat-client")

	once.Do(func() {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatal(fmt.Errorf("cannot make required directory: %w", err))
		}
	})
	return filepath.Join(configDir, tokenFileName)
}
