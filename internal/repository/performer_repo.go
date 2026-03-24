package repository

import (
	"context"
	"database/sql"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
)

// PerformerRepo репозиторий для работы с БД.
type PerformerRepo struct {
	mssql *sql.DB
	logg  *logg.Logger
}

// NewPerformerRepo конструктор репозитория.
func NewPerformerRepo(mssql *sql.DB, logger *logg.Logger) *PerformerRepo {
	return &PerformerRepo{mssql: mssql, logg: logger}
}

type PerformerRepository interface {
	AuthByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (bool, error)
}

// AuthByTabNumAndPass аутентификация по табельному номеру и паролю.
func (p *PerformerRepo) AuthByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (bool, error) {
	var authSuccess bool

	if err := p.mssql.QueryRowContext(ctx, svPerformerAuthQuery, tabNum, passwd).Scan(&authSuccess); err != nil {
		p.logg.LogE(msg.ERS500, err, skipNofS)

		return false, err
	}

	return false, nil
}
