package repository

import (
	"context"

	"gorm.io/gorm"
)

type StatsRepository interface {
	GetReviewAssignmentsCount(ctx context.Context) (map[string]int64, error)
}

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) GetReviewAssignmentsCount(ctx context.Context) (map[string]int64, error) {
	type row struct {
		UserID string
		Cnt    int64
	}

	var rows []row
	err := r.db.WithContext(ctx).
		Table("reviewers").
		Select("user_id, count(*) as cnt").
		Group("user_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string]int64, len(rows))
	for _, row := range rows {
		res[row.UserID] = row.Cnt
	}
	return res, nil
}
