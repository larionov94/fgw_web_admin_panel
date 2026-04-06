package repository

import (
	"context"
	"database/sql"
	"fgw_web_aforms_panel/internal/entity"
	"fgw_web_aforms_panel/pkg/logg"
	"fgw_web_aforms_panel/pkg/msg"
)

// HistoryRepo репозиторий для работы с БД.
type HistoryRepo struct {
	mssql *sql.DB
	logg  *logg.Logger
}

// NewHistoryRepo конструктор репозитория.
func NewHistoryRepo(mssql *sql.DB, logger *logg.Logger) *HistoryRepo {
	return &HistoryRepo{mssql, logger}
}

type HistoryRepository interface {
	AddEntryAndExit(ctx context.Context, performerData *entity.HistoryPerformer) error
}

// AddEntryAndExit Добавление записи в историю входов/выходов.
func (h *HistoryRepo) AddEntryAndExit(ctx context.Context, performerData *entity.HistoryPerformer) error {
	if performerData == nil {
		h.logg.LogW(msg.WRS400, logg.SkipNofS)

		return nil
	}

	if _, err := h.mssql.ExecContext(ctx, svAFHistoryOfEntryAndExitAddQuery,
		performerData.Hostname,
		performerData.IpAddress,
		performerData.TraceId,
		performerData.FIO,
		performerData.RoleName,
		performerData.EntryExit,
		performerData.CreatedBy,
	); err != nil {
		h.logg.LogE(msg.ERS500, err, logg.SkipNofS)

		return err
	}

	return nil
}
