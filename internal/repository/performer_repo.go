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
	AuthByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (*entity.Performer, error)
	FindByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error)
}

// AuthByTabNumAndPass аутентификация по табельному номеру и паролю.
func (p *PerformerRepo) AuthByTabNumAndPass(ctx context.Context, tabNum int, passwd string) (*entity.Performer, error) {
	var performer entity.Performer

	if err := p.mssql.QueryRowContext(ctx, svPerformerAuthQuery, tabNum, passwd).Scan(
		&performer.Id,
		&performer.SectorId,
		&performer.FIO,
		&performer.TabNum,
		&performer.Barcode,
		&performer.AccessBarcode,
		&performer.Archive,
		&performer.PerformerRole.RoleIdAForms,
		&performer.PerformerRole.RoleIdAFGW,
		&performer.PerformerRole.RoleNameAForms,
		&performer.PerformerRole.RoleNameAFGW,
		&performer.AuthSuccess,
	); err != nil {
		p.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return nil, err
	}

	return &performer, nil
}

// FindByTabNum ищет сотрудника по табельному номеру.
func (p *PerformerRepo) FindByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error) {
	var performer entity.Performer

	if err := p.mssql.QueryRowContext(ctx, svPerformerFindByTabNumQuery, tabNum).Scan(
		&performer.FIO,
		&performer.TabNum,
		&performer.PerformerRole.RoleIdAForms,
		&performer.PerformerRole.RoleIdAFGW,
		&performer.PerformerRole.RoleNameAForms,
		&performer.PerformerRole.RoleNameAFGW,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logg.LogEf(logg.SkipNofS, err, "%s: tabNum: %d", msg.ERS501, tabNum)

			return nil, err
		}

		p.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return nil, err
	}

	return &performer, nil
}
