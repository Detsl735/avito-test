package http

import (
	"errors"
	"net/http"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"github.com/Detsl735/avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	prService   service.PRService
	statsRepo   repository.StatsRepository
}

func NewUserHandler(userSvc service.UserService, prSvc service.PRService, statsRepo repository.StatsRepository) *UserHandler {
	return &UserHandler{
		userService: userSvc,
		prService:   prSvc,
		statsRepo:   statsRepo,
	}
}

func (h *UserHandler) Register(r *gin.RouterGroup) {
	r.POST("/users/setIsActive", h.SetIsActive)
	r.GET("/users/getReview", h.GetReview)
	r.GET("/stats", h.Stats) // дополнительный эндпоинт
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest(err.Error()))
		return
	}

	user, err := h.userService.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "user not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}

	c.JSON(http.StatusOK, UserResponse{User: *user})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errorBadRequest("user_id is required"))
		return
	}

	if _, err := h.userService.GetByID(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return

	}

	prs, err := h.prService.GetReviewPRs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}

	resp := GetReviewResponse{
		UserID:       userID,
		PullRequests: prs,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Stats(c *gin.Context) {
	stats, err := h.statsRepo.GetReviewAssignmentsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}
	c.JSON(http.StatusOK, StatsResponse{Assignments: stats})
}
