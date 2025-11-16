package service

import (
	"context"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, active bool) (*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	user, err := s.userRepo.SetIsActive(ctx, userID, active)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}
