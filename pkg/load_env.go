package pkg

import (
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

const formatFileForEnv = ".env"

func LoadEnvFile(pathToFile string) error {
	envPath := filepath.Join(pathToFile, formatFileForEnv)
	err := godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.ES5004, err)
	}

	return nil
}
