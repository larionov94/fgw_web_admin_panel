package pkg

import (
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	formatFileForEnv = ".1env"
	skipNofS         = 4 // skipNofS кол-во пропускаемых кадров стека.
)

func LoadEnvFile(pathToFile string, logger *logg.Logger) error {
	if logger == nil {
		return fmt.Errorf(msg.EL5010)
	}
	envPath := filepath.Join(pathToFile, formatFileForEnv)
	err := godotenv.Load(envPath)
	if err != nil {
		logger.LogEf(skipNofS, err, "%s", msg.ES5004)

		return err
	}

	return nil
}
