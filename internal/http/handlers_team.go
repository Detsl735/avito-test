package http

import (
	"errors"
	"net/http"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) Register(r *gin.RouterGroup) {
	r.POST("/team/add", h.AddTeam)
	r.GET("/team/get", h.GetTeam)
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req TeamAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest(err.Error()))
		return
	}

	team, users, err := h.teamService.AddTeam(c.Request.Context(), req.TeamName, req.Members)
	if err != nil {
		if errors.Is(err, domain.ErrTeamExists) {
			c.JSON(http.StatusBadRequest, errorResponse("TEAM_EXISTS", "team_name already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}

	resp := TeamAddResponse{}
	resp.Team.TeamName = team.TeamName
	for _, u := range users {
		resp.Team.Members = append(resp.Team.Members, domain.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, errorBadRequest("team_name is required"))
		return
	}

	team, users, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "team not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}

	members := make([]domain.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, domain.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"team_name": team.TeamName,
		"members":   members,
	})
}
