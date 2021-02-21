package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	viewercontext "github.com/nicopozo/pr-viewer/internal/context"
	"github.com/nicopozo/pr-viewer/internal/model"
	jsonutils "github.com/nicopozo/pr-viewer/internal/utils/json"
	"github.com/nicopozo/pr-viewer/internal/utils/log"
)

type Client interface {
	GetRepositoryPullRequests(ctx context.Context, owner, repository, token string) (*model.PullRequestList, error)
	GetUsername(ctx context.Context, token string) (string, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) (Client, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("http client can not be nil")
	}

	return &client{httpClient: httpClient}, nil
}

func (cli *client) GetUsername(ctx context.Context, token string) (string, error) {
	logger := viewercontext.Logger(ctx)

	githubURL := "https://api.github.com/graphql"

	q := `query User {viewer {login}}`

	body := query{Query: q}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, githubURL, strings.NewReader(jsonutils.Marshal(body)))
	if err != nil {
		return "", fmt.Errorf("error creating request, %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := cli.httpClient.Do(request) //nolint:bodyclose
	if err != nil {
		err = fmt.Errorf("unable to execute request, %w", err)
		logger.Error(cli, nil, err, "Unable to make request")

		return "", err
	}

	defer closeResponseBody(cli, resp, logger)

	responseBody := ""
	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		responseBody = string(b)
	}

	logger.Debug(cli, nil, "github username response: %s", responseBody)

	response := new(viewerResponse)

	err = jsonutils.Unmarshal(strings.NewReader(responseBody), response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling body, %w", err)
	}

	return response.Data.Viewer.Username, nil
}

func (cli *client) GetRepositoryPullRequests(ctx context.Context, owner, repository,
	token string) (*model.PullRequestList, error) {
	logger := viewercontext.Logger(ctx)

	githubURL := "https://api.github.com/graphql"
	//githubURL := "http://api.mp.internal.ml.com/reconciliations/mockservice/mock/github/graphql"

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
			updatedAt
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
		createdAt
		title
        repository {
          name
        }
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

	response := new(repositoryResponse)

	err = jsonutils.Unmarshal(resp.Body, response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	result := new(model.PullRequestList)

	for _, pr := range response.Data.Repository.PullRequests.PullRequests {
		result.PullRequests = append(result.PullRequests, transformPullRequest(pr))
	}

	result.Total = response.Data.Repository.PullRequests.Total

	return result, nil
}

func transformPullRequest(pr pullRequest) model.PullRequest {
	result := model.PullRequest{
		Author:      pr.Author.Username,
		Assignees:   pr.Assignees.TotalCount,
		Url:         pr.Url,
		CreatedAt:   pr.CreatedAt,
		Application: pr.Repository.Name,
		Title:       pr.Title,
	}

	result.Story = getStory(pr.Title)

	for _, rr := range pr.ReviewRequests.ReviewRequests {
		reviewRequest := model.ReviewRequest{
			RequestedReviewer: rr.RequestedReviewer.Username,
			State:             model.ReviewRequestStatusPending,
		}
		result.ReviewRequests = append(result.ReviewRequests, reviewRequest)
	}

	for _, r := range pr.Reviews.PullRequestReview {
		review := model.PullRequestReview{
			Author:    r.Author.Username,
			State:     r.State,
			UpdatedAt: r.UpdatedAt,
		}

		result.Reviews = append(result.Reviews, review)
	}

	return result
}

func getStory(title string) string {
	story := getStoryFromTitle(title, "MPCON-")
	if story != "" {

		return story
	}

	return getStoryFromTitle(title, "LIQ-")
}

func getStoryFromTitle(title, storyPrefix string) string {
	title = strings.ToUpper(title)
	if strings.HasPrefix(title, storyPrefix) {
		idx := 0
		for pos, char := range title {
			charStr := fmt.Sprintf("%c", char)
			if pos >= len(storyPrefix) {
				if _, err := strconv.Atoi(charStr); err != nil {
					idx = pos
					break
				}
			}
		}

		return title[:idx]
	}

	return ""
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
