package repository

import (
	"context"
	"database/sql"
	"fgw_web_admin_panel/pkg/logg"
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
	FindByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (bool, error)
}

// FindByTabNumAndPass аутентификация по табельному номеру и паролю.
func (p *PerformerRepo) FindByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (bool, error) {

	return false, nil
}
