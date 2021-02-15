package service

import (
	"context"
	"fmt"

	"github.com/nicopozo/pr-viewer/internal/github"
)

type PRService interface {
	GetPRs(ctx context.Context, userType, token string)
}

type githubPRService struct {
	githubClient github.Client
}

func NewGithubPRService(githubClient github.Client) (PRService, error) {
	if githubClient == nil {
		return nil, fmt.Errorf("github client can not be nil")
	}

	return &githubPRService{githubClient: githubClient}, nil
}

func (g githubPRService) GetPRs(ctx context.Context, userType, token string) {
	panic("implement me")
}

/*


	apps := []string{"fury_gateway-reconciliation",
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

	for _, app := range apps {
		result, err := githubClient.GetRepositoryPullRequests(context.Background(), "mercadolibre", app,
			"dc84e45987cca90c2e9e57f802acc85265108796")
		if err != nil {
			println(err.Error())
		} else {
			for _, pr := range result.PullRequests {
				fmt.Println(app, pr.Author, pr.Url, pr.Assignees, pr.ReviewRequests)
			}
		}
	}
*/
