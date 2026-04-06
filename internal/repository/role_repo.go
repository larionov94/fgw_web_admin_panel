package repository

import (
	"context"
	"database/sql"
	"fgw_web_aforms_panel/internal/config/db"
	"fgw_web_aforms_panel/internal/entity"
	"fgw_web_aforms_panel/pkg/logg"
	"fgw_web_aforms_panel/pkg/msg"
)

type RoleRepo struct {
	mssql *sql.DB
	logg  *logg.Logger
}

func NewRoleRepo(mssql *sql.DB, logg *logg.Logger) *RoleRepo {
	return &RoleRepo{mssql: mssql, logg: logg}
}

type RoleRepository interface {
	All(ctx context.Context) ([]*entity.Role, error)
}

// All получает список ролей из БД.
func (r *RoleRepo) All(ctx context.Context) ([]*entity.Role, error) {
	rows, err := r.mssql.QueryContext(ctx, svRolesQuery)
	if err != nil {
		r.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return nil, err
	}

	defer db.CloseRows(rows, r.logg)

	var roles []*entity.Role
	for rows.Next() {
		var role entity.Role
		if err = rows.Scan(
			&role.Id,
			&role.NameRole,
			&role.Description,
			&role.AuditRec.CreatedAt,
			&role.AuditRec.CreatedBy,
			&role.AuditRec.UpdatedAt,
			&role.AuditRec.UpdatedBy,
		); err != nil {
			r.logg.LogE(msg.ERS502, err, logg.SkipNofS)

			return nil, err
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		r.logg.LogE(msg.ERS503, err, logg.SkipNofS)

		return nil, err
	}

	return roles, nil
}
