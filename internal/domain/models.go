package domain

import "time"

type User struct {
	UserID   string `gorm:"column:user_id;primaryKey"`
	Username string `gorm:"column:username;not null"`
	TeamName string `gorm:"column:team_name;not null;index"`
	IsActive bool   `gorm:"column:is_active;not null;default:true"`
}

func (User) TableName() string {
	return "users"
}

type Team struct {
	TeamName string `gorm:"column:team_name;primaryKey"`
}

func (Team) TableName() string {
	return "teams"
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string     `gorm:"column:pull_request_id;primaryKey"`
	PullRequestName string     `gorm:"column:pull_request_name;not null"`
	AuthorID        string     `gorm:"column:author_id;not null;index"`
	Status          PRStatus   `gorm:"column:status;type:varchar(16);not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`
	MergedAt        *time.Time `gorm:"column:merged_at"`
}

func (PullRequest) TableName() string {
	return "pull_requests"
}

type Reviewer struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	PullRequestID string `gorm:"column:pull_request_id;not null;index"`
	UserID        string `gorm:"column:user_id;not null;index"`
}

func (Reviewer) TableName() string {
	return "reviewers"
}

type PullRequestFull struct {
	PullRequest
	AssignedReviewers []string
}

type PullRequestShort struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PRStatus `json:"status"`
}
