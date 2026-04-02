package repository

import (
	"context"
	"database/sql"
	"errors"
	"fgw_web_admin_panel/internal/config/db"
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
	All(ctx context.Context) ([]*entity.Performer, error)
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

// All выводит список сотрудников.
func (p *PerformerRepo) All(ctx context.Context) ([]*entity.Performer, error) {
	rows, err := p.mssql.QueryContext(ctx, svPerformerAllQuery)
	if err != nil {
		p.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return nil, err
	}

	defer db.CloseRows(rows, p.logg)

	var performers []*entity.Performer
	for rows.Next() {
		var performer entity.Performer
		if err = rows.Scan(
			&performer.Id,
			&performer.SectorId,
			&performer.FIO,
			&performer.TabNum,
			&performer.Barcode,
			&performer.AccessBarcode,
			&performer.Passwd,
			&performer.IssuedAt,
			&performer.Archive,
			&performer.PerformerRole.RoleIdAForms,
			&performer.PerformerRole.RoleIdAFGW,
			&performer.AuditRec.CreatedAt,
			&performer.AuditRec.CreatedBy,
			&performer.AuditRec.UpdatedAt,
			&performer.AuditRec.UpdatedBy,
		); err != nil {
			p.logg.LogE(msg.ERS502, err, logg.SkipNofS)

			return nil, err
		}

		performers = append(performers, &performer)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.ERS503, err, logg.SkipNofS)

		return nil, err
	}

	return performers, nil
}
