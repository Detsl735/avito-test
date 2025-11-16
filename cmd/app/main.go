package main

import (
	"log"

	"github.com/Detsl735/avito-test/internal/config"
	"github.com/Detsl735/avito-test/internal/domain"
	transport "github.com/Detsl735/avito-test/internal/http"
	"github.com/Detsl735/avito-test/internal/repository"
	"github.com/Detsl735/avito-test/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := db.AutoMigrate(&domain.Team{}, &domain.User{}, &domain.PullRequest{}, &domain.Reviewer{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	teamSvc := service.NewTeamService(db, teamRepo, userRepo)
	userSvc := service.NewUserService(db, userRepo)
	prSvc := service.NewPRService(db, prRepo, userRepo)

	router := transport.NewRouter(teamSvc, userSvc, prSvc, statsRepo)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
