package github

type Repository struct {
	PullRequests PullRequestList `json:"pullRequests"`
}

type RepositoryResponse struct {
	Data struct {
		Repository Repository `json:"repository"`
	} `json:"data"`
}
