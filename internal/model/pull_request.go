package model

import "time"

type PullRequest struct {
	Author         string              `json:"author"`
	Reviews        []PullRequestReview `json:"reviews"`
	ReviewRequests []ReviewRequest     `json:"review_requests"`
	Assignees      int                 `json:"assignees"`
	Url            string              `json:"url"`
	CreatedAt      time.Time           `json:"created_at"`
	Application    string              `json:"application"`
	State          string              `json:"state"`
	Title          string              `json:"title"`
	Story          string              `json:"story"`
}

//PullRequestList
type PullRequestList struct {
	Total          int             `json:"total"`
	PullRequests   []PullRequest   `json:"pull_requests"`
	ReviewersCount []ReviewerCount `json:"reviewers_count"`
}

type ReviewerCount struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}
