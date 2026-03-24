package config

import (
	"fgw_web_admin_panel/pkg"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"os"
)

type CfgEntryMSSQL struct {
	Server  string `env:"MSSQL_SERVER" `
	Name    string `env:"MSSQL_NAME" `
	User    string `env:"MSSQL_USER" `
	Passwd  string `env:"MSSQL_PASSWD" `
	Charset string `env:"MSSQL_CHARSET" `
}

type CfgMSSQL struct {
	CfgDB  *CfgEntryMSSQL
	logger *logg.Logger
}

// NewCfgMSSQL создает новую конфигурацию MSSQL из .env файла.
func NewCfgMSSQL(logger *logg.Logger) (*CfgMSSQL, error) {
	if logger == nil {
		return nil, fmt.Errorf(msg.EL5010)
	}

	if err := pkg.LoadEnvFile("", logger); err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.ES5004)

		return nil, err
	}

	return &CfgMSSQL{
		&CfgEntryMSSQL{
			Server:  os.Getenv("MSSQL_SERVER"),
			Name:    os.Getenv("MSSQL_NAME"),
			User:    os.Getenv("MSSQL_USER"),
			Passwd:  os.Getenv("MSSQL_PASSWD"),
			Charset: os.Getenv("MSSQL_CHARSET"),
		},
		logger,
	}, nil
}
