package service

import (
	"context"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/internal/repository"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
)

type RoleService struct {
	roleRepo repository.RoleRepository
	logg     *logg.Logger
}

func NewRoleService(roleRepo repository.RoleRepository, logger *logg.Logger) *RoleService {
	return &RoleService{roleRepo: roleRepo, logg: logger}
}

type RoleUseCase interface {
	AllRoles(ctx context.Context) ([]*entity.Role, error)
}

// AllRoles бизнес-логика получает список ролей.
func (r *RoleService) AllRoles(ctx context.Context) ([]*entity.Role, error) {
	roles, err := r.roleRepo.All(ctx)
	if err != nil {
		r.logg.LogE(msg.ESR502, err, logg.SkipNofS)

		return nil, err
	}

	return roles, nil
}
