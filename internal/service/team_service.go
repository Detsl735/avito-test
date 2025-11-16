package service

import (
	"context"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"gorm.io/gorm"
)

type TeamService interface {
	AddTeam(ctx context.Context, teamName string, members []domain.TeamMember) (*domain.Team, []domain.User, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, []domain.User, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewTeamService(db *gorm.DB, tRepo repository.TeamRepository, uRepo repository.UserRepository) TeamService {
	return &teamService{
		teamRepo: tRepo,
		userRepo: uRepo,
		db:       db,
	}
}

func (s *teamService) AddTeam(ctx context.Context, teamName string, members []domain.TeamMember) (*domain.Team, []domain.User, error) {
	if _, err := s.teamRepo.GetByName(ctx, teamName); err == nil {
		return nil, nil, domain.ErrTeamExists
	}

	if err := s.db.WithContext(ctx).Create(&domain.Team{TeamName: teamName}).Error; err != nil {
		return nil, nil, err
	}

	users := make([]domain.User, 0, len(members))
	for _, m := range members {
		users = append(users, domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: teamName,
			IsActive: m.IsActive,
		})
	}

	if err := s.userRepo.UpsertMany(ctx, users); err != nil {
		return nil, nil, err
	}

	return &domain.Team{TeamName: teamName}, users, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, []domain.User, error) {
	t, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}
	users, err := s.userRepo.GetByTeamName(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}
	return t, users, nil
}
