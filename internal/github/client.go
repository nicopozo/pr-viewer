package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	viewercontext "github.com/nicopozo/pr-viewer/internal/context"
	jsonutils "github.com/nicopozo/pr-viewer/internal/utils/json"
	"github.com/nicopozo/pr-viewer/internal/utils/log"
)

type Client interface {
	GetRepositoryPullRequests(ctx context.Context, owner, repository,
		token string) (*PullRequestList, error)
}

type client struct {
	httpClient *http.Client
}

func NewGithubClient(httpClient *http.Client) (Client, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("http client can not be nil")
	}

	return &client{httpClient: httpClient}, nil
}

func (cli *client) GetRepositoryPullRequests(ctx context.Context, owner, repository,
	token string) (*PullRequestList, error) {
	logger := viewercontext.Logger(ctx)

	githubURL := "https://api.github.com/graphql"

	q := `query Repository {
  repository(owner: "%s", name: "%s") {
    pullRequests(states: [OPEN], first: 100) {
      nodes {
        reviews(first: 100) {
          nodes {
            author {
              login
            }
            state
          }
        }
        reviewRequests(first: 100) {
          nodes {
            id
            requestedReviewer {
              ... on User {
                login
              }
            }
          }
        }
        author {
          login
        }
        url
      }
      totalCount
    }
  }
}`

	body := query{Query: fmt.Sprintf(q, owner, repository)}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, githubURL, strings.NewReader(jsonutils.Marshal(body)))
	if err != nil {
		return nil, fmt.Errorf("error creating request, %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := cli.httpClient.Do(request) //nolint:bodyclose
	if err != nil {
		err = fmt.Errorf("unable to make request, %w", err)
		logger.Error(cli, nil, err, "Unable to make request")

		return nil, err
	}

	defer closeResponseBody(cli, resp, logger)

	response := new(RepositoryResponse)

	err = jsonutils.Unmarshal(resp.Body, response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return &response.Data.Repository.PullRequests, nil
}

func closeResponseBody(cli interface{}, response *http.Response, logger log.ILogger) {
	if response != nil && response.Body != nil {
		if err := response.Body.Close(); err != nil {
			errorMsg := "error closing response body"

			if logger != nil {
				logger.Error(cli, nil, err, errorMsg)
			} else {
				fmt.Printf("%v %s\n", reflect.TypeOf(cli), errorMsg)
			}
		}
	}
}

type query struct {
	Query string `json:"query"`
}
