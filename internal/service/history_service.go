package service

import (
	"context"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/internal/repository"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
)

type HistoryService struct {
	historyRepo repository.HistoryRepository
	logg        *logg.Logger
}

func NewHistoryService(historyRepo repository.HistoryRepository, logger *logg.Logger) *HistoryService {
	return &HistoryService{historyRepo: historyRepo, logg: logger}
}

type HistoryUseCase interface {
	AddHistoryOfEntryAndExit(ctx context.Context, performer *entity.HistoryPerformer) error
}

// AddHistoryOfEntryAndExit бизнес-логика ведения истории входа/выхода из программы.
func (h *HistoryService) AddHistoryOfEntryAndExit(ctx context.Context, performerHistory *entity.HistoryPerformer) error {
	if performerHistory == nil {
		h.logg.LogW(msg.WRS400, logg.SkipNofS)

		return nil
	}

	return h.historyRepo.AddEntryAndExit(ctx, performerHistory)
}
