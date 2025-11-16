package repository

import (
	"context"

	"github.com/Detsl735/avito-test/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	UpsertMany(ctx context.Context, users []domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByTeamName(ctx context.Context, teamName string) ([]domain.User, error)
	SetIsActive(ctx context.Context, id string, active bool) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) UpsertMany(ctx context.Context, users []domain.User) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, u := range users {
		var existing domain.User
		err := tx.Where("user_id = ?", u.UserID).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&u).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		} else {
			existing.Username = u.Username
			existing.TeamName = u.TeamName
			existing.IsActive = u.IsActive
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByTeamName(ctx context.Context, teamName string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).Where("team_name = ?", teamName).Find(&users).Error
	return users, err
}

func (r *userRepository) SetIsActive(ctx context.Context, id string, active bool) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	u.IsActive = active
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
