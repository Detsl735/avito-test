package http

import "github.com/Detsl735/avito-test/internal/domain"

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type TeamAddRequest struct {
	TeamName string              `json:"team_name" binding:"required"`
	Members  []domain.TeamMember `json:"members" binding:"required"`
}

type TeamAddResponse struct {
	Team struct {
		TeamName string              `json:"team_name"`
		Members  []domain.TeamMember `json:"members"`
	} `json:"team"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User domain.User `json:"user"`
}

type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type PullRequestResponse struct {
	PR struct {
		PullRequestID   string   `json:"pull_request_id"`
		PullRequestName string   `json:"pull_request_name"`
		AuthorID        string   `json:"author_id"`
		Status          string   `json:"status"`
		Assigned        []string `json:"assigned_reviewers"`
		CreatedAt       string   `json:"createdAt,omitempty"`
		MergedAt        *string  `json:"mergedAt,omitempty"`
	} `json:"pr"`
}

type PullRequestMergeRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type PullRequestReassignRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id" binding:"required"`
}

type PullRequestReassignResponse struct {
	PR struct {
		PullRequestID   string   `json:"pull_request_id"`
		PullRequestName string   `json:"pull_request_name"`
		AuthorID        string   `json:"author_id"`
		Status          string   `json:"status"`
		Assigned        []string `json:"assigned_reviewers"`
		CreatedAt       string   `json:"createdAt,omitempty"`
		MergedAt        *string  `json:"mergedAt,omitempty"`
	} `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

type GetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []domain.PullRequestShort `json:"pull_requests"`
}

type StatsResponse struct {
	Assignments map[string]int64 `json:"assignments"`
}
