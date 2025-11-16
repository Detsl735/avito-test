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

func setupUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Team{}, &domain.User{})
	require.NoError(t, err)

	err = db.Create(&domain.Team{TeamName: "backend"}).Error
	require.NoError(t, err)

	return db
}

func TestUserService_SetIsActive_Success(t *testing.T) {
	db := setupUserTestDB(t)

	userRepo := repository.NewUserRepository(db)
	svc := NewUserService(db, userRepo)

	ctx := context.Background()

	u := domain.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}
	require.NoError(t, db.Create(&u).Error)

	updated, err := svc.SetIsActive(ctx, "u1", false)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "u1", updated.UserID)
	require.False(t, updated.IsActive)

	var dbUser domain.User
	err = db.First(&dbUser, "user_id = ?", "u1").Error
	require.NoError(t, err)
	require.False(t, dbUser.IsActive)
}

func TestUserService_SetIsActive_NotFound(t *testing.T) {
	db := setupUserTestDB(t)

	userRepo := repository.NewUserRepository(db)
	svc := NewUserService(db, userRepo)

	ctx := context.Background()

	user, err := svc.SetIsActive(ctx, "no-such-user", true)
	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, domain.ErrNotFound, err)
}

func TestUserService_GetByID_Success(t *testing.T) {
	db := setupUserTestDB(t)

	userRepo := repository.NewUserRepository(db)
	svc := NewUserService(db, userRepo)

	u := domain.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}
	require.NoError(t, db.Create(&u).Error)

	ctx := context.Background()

	got, err := svc.GetByID(ctx, "u1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "u1", got.UserID)
	require.Equal(t, "Alice", got.Username)
	require.Equal(t, "backend", got.TeamName)
	require.True(t, got.IsActive)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	db := setupUserTestDB(t)

	userRepo := repository.NewUserRepository(db)
	svc := NewUserService(db, userRepo)

	ctx := context.Background()

	got, err := svc.GetByID(ctx, "no-such-user")
	require.Error(t, err)
	require.Nil(t, got)
	require.Equal(t, domain.ErrNotFound, err)
}
