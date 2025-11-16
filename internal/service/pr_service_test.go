package service

import (
	"context"
	"testing"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"github.com/stretchr/testify/require"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Team{}, &domain.User{}, &domain.PullRequest{}, &domain.Reviewer{})
	require.NoError(t, err)

	return db
}

func TestCreatePR_AssignsReviewers(t *testing.T) {
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)
	prSvc := NewPRService(db, prRepo, userRepo)

	ctx := context.Background()

	err := db.Create(&domain.Team{TeamName: "backend"}).Error
	require.NoError(t, err)

	users := []domain.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true},
	}
	require.NoError(t, userRepo.UpsertMany(ctx, users))

	pr, err := prSvc.CreatePR(ctx, "pr-1", "Test", "u1")
	require.NoError(t, err)
	require.Equal(t, "pr-1", pr.PullRequestID)
	require.Equal(t, domain.PRStatusOpen, pr.Status)
	require.True(t, len(pr.AssignedReviewers) <= 2)
	for _, r := range pr.AssignedReviewers {
		require.NotEqual(t, "u1", r)
	}
}

func TestMergePR_Idempotent(t *testing.T) {
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)
	prSvc := NewPRService(db, prRepo, userRepo)

	ctx := context.Background()

	err := db.Create(&domain.PullRequest{
		PullRequestID:   "pr-2",
		PullRequestName: "Test merge",
		AuthorID:        "u1",
		Status:          domain.PRStatusOpen,
	}).Error
	require.NoError(t, err)

	full, err := prSvc.MergePR(ctx, "pr-2")
	require.NoError(t, err)
	require.Equal(t, domain.PRStatusMerged, full.Status)
	require.NotNil(t, full.MergedAt)

	full2, err := prSvc.MergePR(ctx, "pr-2")
	require.NoError(t, err)
	require.Equal(t, domain.PRStatusMerged, full2.Status)
}
