package service

import (
	"context"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/internal/repository"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
)

type PerformerService struct {
	performerRepo repository.PerformerRepository
	logg          *logg.Logger
}

func NewPerformerUseCase(performerRepo repository.PerformerRepository) *PerformerService {
	return &PerformerService{performerRepo: performerRepo}
}

type PerformerUseCase interface {
	AuthPerformer(ctx context.Context, tabNum int, passwd string) (string, error)
}

// AuthPerformer бизнес-логика аутентификации сотрудника.
func (p *PerformerService) AuthPerformer(ctx context.Context, tabNum int, passwd string) (*entity.PerformerAuth, error) {
	if tabNum <= 0 || passwd == "" {
		p.logg.LogW(msg.WSR400, logg.SkipNofS)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.WSR400,
		}, nil
	}

	authOK, err := p.performerRepo.AuthByTabNumAndPass(ctx, tabNum, passwd)
	if err != nil {
		p.logg.LogE(msg.ESR501, err, logg.SkipNofS)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.ESR501,
		}, nil
	}

	if !authOK {
		p.logg.LogWf(logg.SkipNofS, "%s: Таб. номер: %d", msg.WSR401, tabNum)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.WSR401,
		}, err
	}

	return &entity.PerformerAuth{
		Success: true,
		Message: msg.ISR200,
	}, nil
}
