package domain

import "errors"

var (
	ErrTeamExists = errors.New("team already exists")

	ErrPRExists    = errors.New("pr already exists")
	ErrPRMerged    = errors.New("pr already merged")
	ErrNotAssigned = errors.New("user is not assigned as reviewer")
	ErrNoCandidate = errors.New("no candidate for reviewer")
	ErrNotFound    = errors.New("not found")
)
