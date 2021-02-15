package github

import "time"

type PullRequest struct {
	Author         User                  `json:"author"`
	Reviews        PullRequestReviewList `json:"reviews"`
	ReviewRequests ReviewRequestList     `json:"reviewRequests"`
	Assignees      Assignees             `json:"assignees"`
	Url            string                `json:"url"`
	CreatedAt      time.Time             `json:"createdAt"`
}

//PullRequestList represents a Github PullRequestConnection object
type PullRequestList struct {
	PullRequests []PullRequest `json:"nodes"`
	Total        int           `json:"totalCount"`
}

type Assignees struct {
	TotalCount int `json:"totalCount"`
}
