package repository

import (
	"context"
	"database/sql"
	"errors"
	"fgw_web_admin_panel/internal/entity"
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
	FindByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error)
}

// AuthByTabNumAndPass аутентификация по табельному номеру и паролю.
func (p *PerformerRepo) AuthByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (bool, error) {
	var authSuccess bool

	if err := p.mssql.QueryRowContext(ctx, svPerformerAuthQuery, tabNum, passwd).Scan(&authSuccess); err != nil {
		p.logg.LogE(msg.ERS500, err, skipNofS)

		return false, err
	}

	return authSuccess, nil
}

// FindByTabNum ищет сотрудника по табельному номеру.
func (p *PerformerRepo) FindByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error) {
	var performer entity.Performer

	if err := p.mssql.QueryRowContext(ctx, svPerformerFindByTabNumQuery, tabNum).Scan(
		&performer.Id,
		&performer.SectorId,
		&performer.FIO,
		&performer.TabNum,
		&performer.Barcode,
		&performer.AccessBarcode,
		&performer.Passwd,
		&performer.IssuedAt,
		&performer.Archive,
		&performer.RoleIdAForms,
		&performer.RoleIdAFGW,
		&performer.AuditRec.CreatedAt,
		&performer.AuditRec.CreatedBy,
		&performer.AuditRec.UpdatedAt,
		&performer.AuditRec.UpdatedBy,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logg.LogEf(skipNofS, err, "%s: tabNum: %d", msg.ERS501, tabNum)

			return nil, err
		}

		p.logg.LogE(msg.ERS500, err, skipNofS)

		return nil, err
	}

	return &performer, nil
}
