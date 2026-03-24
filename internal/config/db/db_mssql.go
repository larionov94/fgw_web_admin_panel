package db

import (
	"context"
	"database/sql"
	"fgw_web_admin_panel/internal/config"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	maxOpenConn     = 50
	maxIdleConn     = 50
	maxConnLifetime = 60 * time.Minute
	maxConnIdleTime = 10 * time.Second
)

// NewConnMSSQL создает новое подключение к БД MSSQL.
func NewConnMSSQL(ctx context.Context, cfgMSSQL *config.CfgMSSQL, logger *logg.Logger) (*sql.DB, error) {
	if cfgMSSQL == nil {
		return nil, fmt.Errorf(msg.EDB506)
	}

	if logger == nil {
		return nil, fmt.Errorf(msg.EL5010)
	}

	dataSourceName := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&charset=%s",
		cfgMSSQL.CfgDB.User,
		cfgMSSQL.CfgDB.Passwd,
		cfgMSSQL.CfgDB.Server,
		cfgMSSQL.CfgDB.Name,
		cfgMSSQL.CfgDB.Charset,
	)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB500)

		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetConnMaxIdleTime(maxConnIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		Close(db, logger)
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB501)

		return nil, err
	}
	logger.LogIf(logg.SkipNofS, "%s", msg.IDB200)

	return db, nil
}

// Close закрывает соединение с БД.
func Close(db *sql.DB, logger *logg.Logger) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB502)

		return
	}

	logger.LogIf(logg.SkipNofS, "%s", msg.IDB201)
}

// CloseRows закрывает строки.
func CloseRows(rows *sql.Rows, logger *logg.Logger) {
	if rows == nil {
		return
	}

	if err := rows.Close(); err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB503)

		return
	}

	logger.LogIf(logg.SkipNofS, "%s", msg.IDB202)
}
