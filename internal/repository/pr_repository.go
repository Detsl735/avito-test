package repository

import (
	"context"

	"github.com/Detsl735/avito-test/internal/domain"
	"gorm.io/gorm"
)

type PRRepository interface {
	Create(ctx context.Context, pr domain.PullRequest, reviewers []string) (*domain.PullRequestFull, error)
	GetByID(ctx context.Context, id string) (*domain.PullRequestFull, error)
	Update(ctx context.Context, pr domain.PullRequest, reviewers []string) (*domain.PullRequestFull, error)
	GetByReviewer(ctx context.Context, userID string) ([]domain.PullRequestShort, error)
}

type prRepository struct {
	db *gorm.DB
}

func NewPRRepository(db *gorm.DB) PRRepository {
	return &prRepository{db: db}
}

func (r *prRepository) Create(ctx context.Context, pr domain.PullRequest, reviewers []string) (*domain.PullRequestFull, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(&pr).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, uid := range reviewers {
		if err := tx.Create(&domain.Reviewer{
			PullRequestID: pr.PullRequestID,
			UserID:        uid,
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &domain.PullRequestFull{
		PullRequest:       pr,
		AssignedReviewers: reviewers,
	}, nil
}

func (r *prRepository) GetByID(ctx context.Context, id string) (*domain.PullRequestFull, error) {
	var pr domain.PullRequest
	if err := r.db.WithContext(ctx).First(&pr, "pull_request_id = ?", id).Error; err != nil {
		return nil, err
	}

	var reviewers []domain.Reviewer
	if err := r.db.WithContext(ctx).Where("pull_request_id = ?", id).Find(&reviewers).Error; err != nil {
		return nil, err
	}

	res := &domain.PullRequestFull{PullRequest: pr}
	for _, rv := range reviewers {
		res.AssignedReviewers = append(res.AssignedReviewers, rv.UserID)
	}
	return res, nil
}

func (r *prRepository) Update(ctx context.Context, pr domain.PullRequest, reviewers []string) (*domain.PullRequestFull, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Save(&pr).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if reviewers != nil {
		if err := tx.Where("pull_request_id = ?", pr.PullRequestID).Delete(&domain.Reviewer{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, uid := range reviewers {
			if err := tx.Create(&domain.Reviewer{
				PullRequestID: pr.PullRequestID,
				UserID:        uid,
			}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &domain.PullRequestFull{
		PullRequest:       pr,
		AssignedReviewers: reviewers,
	}, nil
}

func (r *prRepository) GetByReviewer(ctx context.Context, userID string) ([]domain.PullRequestShort, error) {
	var rows []struct {
		PullRequestID   string
		PullRequestName string
		AuthorID        string
		Status          string
	}
	err := r.db.WithContext(ctx).Table("pull_requests pr").
		Select("pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status").
		Joins("JOIN reviewers r ON r.pull_request_id = pr.pull_request_id").
		Where("r.user_id = ?", userID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.PullRequestShort, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.PullRequestShort{
			PullRequestID:   row.PullRequestID,
			PullRequestName: row.PullRequestName,
			AuthorID:        row.AuthorID,
			Status:          domain.PRStatus(row.Status),
		})
	}
	return result, nil
}
