package repository

import (
	"context"

	"github.com/Detsl735/avito-test/internal/domain"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) error
	GetByName(ctx context.Context, teamName string) (*domain.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team domain.Team) error {
	return r.db.WithContext(ctx).Create(&team).Error
}

func (r *teamRepository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	var t domain.Team
	if err := r.db.WithContext(ctx).First(&t, "team_name = ?", teamName).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
