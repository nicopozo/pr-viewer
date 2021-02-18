package github

type repository struct {
	Name         string          `json:"name"`
	PullRequests pullRequestList `json:"pullRequests"`
}

type repositoryResponse struct {
	Data struct {
		Repository repository `json:"repository"`
	} `json:"data"`
}
