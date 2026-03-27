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

func NewPerformerService(performerRepo repository.PerformerRepository, logg *logg.Logger) *PerformerService {
	return &PerformerService{performerRepo: performerRepo, logg: logg}
}

type PerformerUseCase interface {
	AuthPerformerWithData(ctx context.Context, tabNum int, passwd string) (*entity.PerformerAuth, error)
	FindPerformerByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error)
}

// AuthPerformerWithData бизнес-логика аутентификации сотрудника с данными.
func (p *PerformerService) AuthPerformerWithData(ctx context.Context, tabNum int, passwd string) (*entity.PerformerAuth, error) {
	if tabNum <= 0 || passwd == "" {
		p.logg.LogW(msg.WSR400, logg.SkipNofS)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.WSR400,
		}, nil
	}

	authWithDataPerformer, err := p.performerRepo.AuthByTabNumAndPass(ctx, tabNum, passwd)
	if err != nil {
		p.logg.LogE(msg.ESR501, err, logg.SkipNofS)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.ESR501,
		}, err
	}

	if !authWithDataPerformer.AuthSuccess {
		p.logg.LogWf(logg.SkipNofS, "%s: %d", msg.WSR401, tabNum)

		return &entity.PerformerAuth{
			Success: false,
			Message: msg.WSR401,
		}, nil
	}

	return &entity.PerformerAuth{
		Success:   true,
		Performer: *authWithDataPerformer,
		Message:   msg.ISR200,
	}, nil
}

func (p *PerformerService) FindPerformerByTabNum(ctx context.Context, tabNum int) (*entity.Performer, error) {
	performer, err := p.performerRepo.FindByTabNum(ctx, tabNum)
	if err != nil {
		p.logg.LogE(msg.ESR500, err, logg.SkipNofS)

		return nil, err
	}

	return performer, nil
}
