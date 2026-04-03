package repository

import (
	"context"
	"database/sql"
	"fgw_web_admin_panel/internal/config/db"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
)

type SectorRepo struct {
	mssql *sql.DB
	logg  *logg.Logger
}

func NewSectorRepo(mssql *sql.DB, logger *logg.Logger) *SectorRepo {
	return &SectorRepo{mssql: mssql, logg: logger}
}

type SectorRepository interface {
	All(ctx context.Context) ([]*entity.Sector, error)
}

// All получает список печей из БД.
func (s *SectorRepo) All(ctx context.Context) ([]*entity.Sector, error) {
	rows, err := s.mssql.QueryContext(ctx, svAFSectorsAllQuery)
	if err != nil {
		s.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return nil, err
	}

	defer db.CloseRows(rows, s.logg)

	var sectors []*entity.Sector
	for rows.Next() {
		var sector entity.Sector
		if err = rows.Scan(
			&sector.Id,
			&sector.NameSector,
			&sector.VpMlSector,
			&sector.AuditRec.CreatedAt,
			&sector.AuditRec.CreatedBy,
			&sector.AuditRec.UpdatedAt,
			&sector.AuditRec.UpdatedBy,
		); err != nil {
			s.logg.LogE(msg.ERS502, err, logg.SkipNofS)

			return nil, err
		}

		sectors = append(sectors, &sector)
	}

	if err = rows.Err(); err != nil {
		s.logg.LogE(msg.ERS503, err, logg.SkipNofS)

		return nil, err
	}

	return sectors, nil
}
