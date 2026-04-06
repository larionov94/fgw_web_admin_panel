package service

import (
	"context"
	"fgw_web_aforms_panel/internal/entity"
	"fgw_web_aforms_panel/internal/repository"
	"fgw_web_aforms_panel/pkg/logg"
	"fgw_web_aforms_panel/pkg/msg"
)

type SectorService struct {
	sectorRepo repository.SectorRepository
	logg       *logg.Logger
}

func NewSectorService(repo repository.SectorRepository, logger *logg.Logger) *SectorService {
	return &SectorService{repo, logger}
}

type SectorUseCase interface {
	AllSectors(ctx context.Context) ([]*entity.Sector, error)
}

// AllSectors бизнес-логика получает список участков печей.
func (s *SectorService) AllSectors(ctx context.Context) ([]*entity.Sector, error) {
	sectors, err := s.sectorRepo.All(ctx)
	if err != nil {
		s.logg.LogE(msg.ESR502, err, logg.SkipNofS)

		return nil, err
	}

	return sectors, nil
}
