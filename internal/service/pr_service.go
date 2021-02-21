package service

import (
	"context"
	"fmt"
	"sort"
	"sync"

	viewercontext "github.com/nicopozo/pr-viewer/internal/context"
	"github.com/nicopozo/pr-viewer/internal/github"
	"github.com/nicopozo/pr-viewer/internal/model"
	"golang.org/x/sync/errgroup"
)

type PRService interface {
	GetUser(ctx context.Context, token string) (*model.User, error)
	GetPRs(ctx context.Context, userType, token string) (*model.PullRequestList, error)
}

type githubPRService struct {
	githubClient github.Client
	mx           sync.Mutex
	incrMx       sync.Mutex
}

func NewGithubPRService(githubClient github.Client) (PRService, error) {
	if githubClient == nil {
		return nil, fmt.Errorf("github client can not be nil")
	}

	return &githubPRService{githubClient: githubClient}, nil
}

func (svc *githubPRService) GetUser(ctx context.Context, token string) (*model.User, error) {
	logger := viewercontext.Logger(ctx)

	username, err := svc.githubClient.GetUsername(context.Background(), token)
	if err != nil {
		logger.Error(svc, nil, err, "error getting username")

		return nil, fmt.Errorf("error getting username, %w", err)
	}

	return &model.User{
		Username: username,
	}, nil
}

func (svc *githubPRService) GetPRs(ctx context.Context, userType, token string) (*model.PullRequestList, error) {
	logger := viewercontext.Logger(ctx)

	apps := getApps()

	username, err := svc.githubClient.GetUsername(context.Background(), token)
	if err != nil {
		logger.Error(svc, nil, err, "error getting username")

		return nil, fmt.Errorf("error getting username, %w", err)
	}

	result := new(model.PullRequestList)
	reviewersCount := make(map[string]int)
	routinesLimit := make(chan string, 10)
	defer close(routinesLimit)

	var group errgroup.Group

	for i := range apps {
		app := apps[i]
		routinesLimit <- app

		group.Go(func() error {
			resp, err := svc.githubClient.GetRepositoryPullRequests(context.Background(), "mercadolibre", app, token)
			if err != nil {
				logger.Error(svc, nil, err, "error getting PRs for app: %s", app)

				return fmt.Errorf("error getting PRs for app: %s, %w", app, err)
			}

			for _, pr := range resp.PullRequests {
				resultPR := completeReviewer(pr)

				if prApplies(pr, userType, username) {
					svc.mx.Lock()
					result.PullRequests = append(result.PullRequests, resultPR)
					svc.mx.Unlock()
				}

				for _, rr := range resultPR.ReviewRequests {
					if rr.RequestedReviewer != "rp-workflow" {
						svc.incrMx.Lock()
						reviewersCount[rr.RequestedReviewer]++
						svc.incrMx.Unlock()
					}
				}
			}

			<-routinesLimit

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		logger.Error(svc, nil, err, "githubPRService.GetPRs() complete with error")
	}

	for rev, c := range reviewersCount {
		count := model.ReviewerCount{
			Username: rev,
			Count:    c,
		}
		result.ReviewersCount = append(result.ReviewersCount, count)
	}

	sort.Slice(result.ReviewersCount, func(i, j int) bool {
		return result.ReviewersCount[i].Count > result.ReviewersCount[j].Count
	})

	result.Total = len(result.PullRequests)

	return result, nil
}

func completeReviewer(pr model.PullRequest) model.PullRequest {
	states := make(map[string]model.PullRequestReview)

	for i := range pr.Reviews {
		newReview := pr.Reviews[i]

		if newReview.Author != pr.Author {
			current, ok := states[newReview.Author]
			if !ok {
				states[newReview.Author] = newReview
			} else {
				switch newReview.State {
				case model.ReviewRequestStatusApproved, model.ReviewRequestStatusChangesRequested:
					if current.State != model.ReviewRequestStatusCommented &&
						current.UpdatedAt.Before(newReview.UpdatedAt) {
						states[newReview.Author] = newReview
					}
					if current.State == model.ReviewRequestStatusApproved &&
						current.UpdatedAt.Before(newReview.UpdatedAt) {
						states[newReview.Author] = newReview
					}
					if current.State == model.ReviewRequestStatusDismissed &&
						current.UpdatedAt.Before(newReview.UpdatedAt) {
						states[newReview.Author] = newReview
					}
				case model.ReviewRequestStatusDismissed:
					if current.State == model.ReviewRequestStatusCommented &&
						current.UpdatedAt.Before(newReview.UpdatedAt) {
						states[newReview.Author] = newReview
					}
				}
			}
		}
	}

	for key, value := range states {
		found := false
		for _, rr := range pr.ReviewRequests {
			if rr.RequestedReviewer == key {
				found = true

				break
			}
		}

		if !found {
			pr.ReviewRequests = append(pr.ReviewRequests, model.ReviewRequest{
				RequestedReviewer: key,
				State:             value.State,
			})
		}
	}

	counts := make(map[string]int)
	counts[model.ReviewRequestStatusPending] = 0
	counts[model.ReviewRequestStatusChangesRequested] = 0
	counts[model.ReviewRequestStatusApproved] = 0
	counts[model.ReviewRequestStatusCommented] = 0
	counts[model.ReviewRequestStatusDismissed] = 0

	for _, rr := range pr.ReviewRequests {
		counts[rr.State] = counts[rr.State] + 1
	}

	if counts[model.ReviewRequestStatusChangesRequested] > 0 {
		pr.State = model.ReviewRequestStatusChangesRequested
	} else if counts[model.ReviewRequestStatusApproved] >= 2 {
		pr.State = model.ReviewRequestStatusApproved
	} else if counts[model.ReviewRequestStatusApproved] == 1 && len(pr.ReviewRequests) == 1 {
		pr.State = model.ReviewRequestStatusApproved
	} else if counts[model.ReviewRequestStatusCommented] > 0 {
		pr.State = model.ReviewRequestStatusCommented
	} else {
		pr.State = model.ReviewRequestStatusPending
	}

	return pr
}

func prApplies(pullRequest model.PullRequest, userType, username string) bool {
	if (userType == "owner" || userType == "") && username == pullRequest.Author {
		return true
	}

	if userType != "owner" {
		for _, reviewRequest := range pullRequest.ReviewRequests {
			if username == reviewRequest.RequestedReviewer {
				return true
			}
		}

		for _, review := range pullRequest.Reviews {
			if username == review.Author {
				return true
			}
		}
	}

	return false
}

func getApps() []string {
	return []string{
		"fury_gateway-reconciliation",
		"fury_mpcs-reconciliation-etl",
		"fury_mpcs-reconciliation-etl-mo",
		"fury_post-recon-outputs",
		"fury_recon-cards-etl",
		"fury_recon-clients-go",
		"fury_recon-commons-go",
		"fury_recon-core",
		"fury_recon-debit-card-etl",
		"fury_recon-drain-tool",
		"fury_recon-file-transfer",
		"fury_recon-forwarder",
		"fury_recon-mocks",
		"fury_recon-online-payments-etl",
		"fury_recon-report-manager",
		"fury_recon-router",
		"fury_recon-simetrik-gateway",
		"fury_recon-sp-prepaid-etl",
		"fury_recon-sp-recharges-etl",
		"fury_recon-sp-utilities-etl",
		"fury_recon-test-utils-go",
		"fury_recon-tools",
		"fury_recon-xoom-report",
		"fury_sodexo-conciliation",
	}
}
