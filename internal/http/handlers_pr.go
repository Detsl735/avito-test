package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/Detsl735/avito-test/internal/domain"

	"github.com/Detsl735/avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

type PRHandler struct {
	prService service.PRService
}

func NewPRHandler(prSvc service.PRService) *PRHandler {
	return &PRHandler{prService: prSvc}
}

func (h *PRHandler) Register(r *gin.RouterGroup) {
	r.POST("/pullRequest/create", h.CreatePR)
	r.POST("/pullRequest/merge", h.MergePR)
	r.POST("/pullRequest/reassign", h.Reassign)
}

func (h *PRHandler) CreatePR(c *gin.Context) {
	var req PullRequestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest(err.Error()))
		return
	}

	full, err := h.prService.CreatePR(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPRExists):
			c.JSON(http.StatusConflict, errorResponse("PR_EXISTS", "PR id already exists"))
			return
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "author or team not found"))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
			return
		}
	}

	c.JSON(http.StatusCreated, prToResponse(full))
}

func (h *PRHandler) MergePR(c *gin.Context) {
	var req PullRequestMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest(err.Error()))
		return
	}

	full, err := h.prService.MergePR(c.Request.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "pr not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
		return
	}

	c.JSON(http.StatusOK, prToResponse(full))
}

func (h *PRHandler) Reassign(c *gin.Context) {
	var req PullRequestReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest(err.Error()))
		return
	}

	full, replacedBy, err := h.prService.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "pr or user not found"))
			return
		case errors.Is(err, domain.ErrPRMerged):
			c.JSON(http.StatusConflict, errorResponse("PR_MERGED", "cannot reassign on merged PR"))
			return
		case errors.Is(err, domain.ErrNotAssigned):
			c.JSON(http.StatusConflict, errorResponse("NOT_ASSIGNED", "reviewer is not assigned to this PR"))
			return
		case errors.Is(err, domain.ErrNoCandidate):
			c.JSON(http.StatusConflict, errorResponse("NO_CANDIDATE", "no active replacement candidate in team"))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("INTERNAL", err.Error()))
			return
		}
	}

	resp := PullRequestReassignResponse{
		ReplacedBy: replacedBy,
	}
	resp.PR = prToResponse(full).PR

	c.JSON(http.StatusOK, resp)
}

func prToResponse(full *domain.PullRequestFull) PullRequestResponse {
	resp := PullRequestResponse{}
	resp.PR.PullRequestID = full.PullRequestID
	resp.PR.PullRequestName = full.PullRequestName
	resp.PR.AuthorID = full.AuthorID
	resp.PR.Status = string(full.Status)
	resp.PR.Assigned = full.AssignedReviewers
	if !full.CreatedAt.IsZero() {
		resp.PR.CreatedAt = full.CreatedAt.UTC().Format(time.RFC3339)
	}
	if full.MergedAt != nil {
		t := full.MergedAt.UTC().Format(time.RFC3339)
		resp.PR.MergedAt = &t
	}
	return resp
}
