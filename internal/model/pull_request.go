package model

import "time"

type PullRequest struct {
	Author         string              `json:"author"`
	Reviews        []PullRequestReview `json:"reviews"`
	ReviewRequests []ReviewRequest     `json:"review_requests"`
	Assignees      int                 `json:"assignees"`
	Url            string              `json:"url"`
	CreatedAt      time.Time           `json:"created_at"`
}

//PullRequestList
type PullRequestList struct {
	Total        int           `json:"total"`
	PullRequests []PullRequest `json:"pull_requests"`
}
