package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/Detsl735/avito-test/internal/domain"
	"github.com/Detsl735/avito-test/internal/repository"
	"gorm.io/gorm"
)

type PRService interface {
	CreatePR(ctx context.Context, id, name, authorID string) (*domain.PullRequestFull, error)
	MergePR(ctx context.Context, id string) (*domain.PullRequestFull, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequestFull, string, error)
	GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequestShort, error)
}

type prService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewPRService(db *gorm.DB, prRepo repository.PRRepository, userRepo repository.UserRepository) PRService {
	return &prService{
		prRepo:   prRepo,
		userRepo: userRepo,
		db:       db,
	}

}

func (s *prService) CreatePR(ctx context.Context, id, name, authorID string) (*domain.PullRequestFull, error) {
	_, err := s.prRepo.GetByID(ctx, id)
	if err == nil {
		return nil, domain.ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, authorID)

	users, err := s.userRepo.GetByTeamName(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	candidates := make([]string, 0)
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if u.UserID == author.UserID {
			continue
		}
		candidates = append(candidates, u.UserID)
	}

	assigned := pickRandom(candidates, 2)

	pr := domain.PullRequest{
		PullRequestID:   id,
		PullRequestName: name,
		AuthorID:        authorID,
		Status:          domain.PRStatusOpen,
		CreatedAt:       time.Now().UTC(),
	}

	return s.prRepo.Create(ctx, pr, assigned)
}

func (s *prService) MergePR(ctx context.Context, id string) (*domain.PullRequestFull, error) {
	full, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if full.Status == domain.PRStatusMerged {
		return full, nil
	}

	now := time.Now().UTC()
	full.Status = domain.PRStatusMerged
	full.MergedAt = &now

	updated, err := s.prRepo.Update(ctx, full.PullRequest, full.AssignedReviewers)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequestFull, string, error) {
	full, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", err
	}

	if full.Status == domain.PRStatusMerged {
		return nil, "", domain.ErrPRMerged
	}

	oldUser, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", err
	}

	found := false
	for _, r := range full.AssignedReviewers {
		if r == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", domain.ErrNotAssigned
	}

	users, err := s.userRepo.GetByTeamName(ctx, oldUser.TeamName)
	if err != nil {
		return nil, "", err
	}

	var candidates []string
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if u.UserID == full.AuthorID {
			continue
		}
		skip := false
		for _, assigned := range full.AssignedReviewers {
			if assigned == u.UserID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		candidates = append(candidates, u.UserID)
	}

	if len(candidates) == 0 {
		return nil, "", domain.ErrNoCandidate
	}

	newUserID := candidates[rand.Intn(len(candidates))]

	for i, rID := range full.AssignedReviewers {
		if rID == oldUserID {
			full.AssignedReviewers[i] = newUserID
			break
		}
	}

	updated, err := s.prRepo.Update(ctx, full.PullRequest, full.AssignedReviewers)
	if err != nil {
		return nil, "", err
	}
	return updated, newUserID, nil
}

func (s *prService) GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequestShort, error) {
	return s.prRepo.GetByReviewer(ctx, userID)
}

func pickRandom(items []string, n int) []string {
	if len(items) == 0 || n <= 0 {
		return nil
	}
	if len(items) <= n {
		return items
	}
	res := make([]string, len(items))
	copy(res, items)
	for i := range res {
		j := rand.Intn(i + 1)
		res[i], res[j] = res[j], res[i]
	}
	return res[:n]
}
