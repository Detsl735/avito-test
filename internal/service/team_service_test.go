package service

import (
	"context"
	"testing"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTeamTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Team{}, &domain.User{})
	require.NoError(t, err)

	return db
}

func TestTeamService_AddTeam_Success(t *testing.T) {
	db := setupTeamTestDB(t)

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewTeamService(db, teamRepo, userRepo)

	ctx := context.Background()

	members := []domain.TeamMember{
		{UserID: "u1", Username: "Alice", IsActive: true},
		{UserID: "u2", Username: "Bob", IsActive: false},
	}

	team, users, err := svc.AddTeam(ctx, "backend", members)
	require.NoError(t, err)
	require.NotNil(t, team)
	require.Equal(t, "backend", team.TeamName)

	require.Len(t, users, 2)
	require.Equal(t, "u1", users[0].UserID)
	require.Equal(t, "backend", users[0].TeamName)
	require.Equal(t, "u2", users[1].UserID)

	var dbTeam domain.Team
	err = db.First(&dbTeam, "team_name = ?", "backend").Error
	require.NoError(t, err)

	var dbUsers []domain.User
	err = db.Where("team_name = ?", "backend").Order("user_id").Find(&dbUsers).Error
	require.NoError(t, err)
	require.Len(t, dbUsers, 2)
}

func TestTeamService_AddTeam_TeamExists(t *testing.T) {
	db := setupTeamTestDB(t)

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewTeamService(db, teamRepo, userRepo)

	ctx := context.Background()

	err := db.Create(&domain.Team{TeamName: "backend"}).Error
	require.NoError(t, err)

	_, _, err = svc.AddTeam(ctx, "backend", []domain.TeamMember{
		{UserID: "u1", Username: "Alice", IsActive: true},
	})
	require.Error(t, err)
	require.Equal(t, domain.ErrTeamExists, err)
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	db := setupTeamTestDB(t)

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewTeamService(db, teamRepo, userRepo)

	err := db.Create(&domain.Team{TeamName: "backend"}).Error
	require.NoError(t, err)

	users := []domain.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: false},
	}
	err = db.Create(&users).Error
	require.NoError(t, err)

	ctx := context.Background()

	team, gotUsers, err := svc.GetTeam(ctx, "backend")
	require.NoError(t, err)
	require.NotNil(t, team)
	require.Equal(t, "backend", team.TeamName)

	require.Len(t, gotUsers, 2)
	ids := []string{gotUsers[0].UserID, gotUsers[1].UserID}
	require.ElementsMatch(t, []string{"u1", "u2"}, ids)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	db := setupTeamTestDB(t)

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewTeamService(db, teamRepo, userRepo)

	ctx := context.Background()

	team, users, err := svc.GetTeam(ctx, "no-such-team")
	require.Error(t, err)
	require.Nil(t, team)
	require.Nil(t, users)
}
