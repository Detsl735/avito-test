package http

import (
	"net/http"

	"github.com/Detsl735/avito-test/internal/repository"
	"github.com/Detsl735/avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	teamSvc service.TeamService,
	userSvc service.UserService,
	prSvc service.PRService,
	statsRepo repository.StatsRepository,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/")
	{
		NewTeamHandler(teamSvc).Register(api)
		NewUserHandler(userSvc, prSvc, statsRepo).Register(api)
		NewPRHandler(prSvc).Register(api)
	}

	return r
}
