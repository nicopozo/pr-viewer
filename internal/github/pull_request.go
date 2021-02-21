package github

import "time"

type pullRequest struct {
	Author         user                  `json:"author"`
	Reviews        pullRequestReviewList `json:"reviews"`
	ReviewRequests reviewRequestList     `json:"reviewRequests"`
	Assignees      assignees             `json:"assignees"`
	Url            string                `json:"url"`
	CreatedAt      time.Time             `json:"createdAt"`
	Repository     repository            `json:"repository"`
	Title          string                `json:"title"`
}

//pullRequestList represents a Github PullRequestConnection object
type pullRequestList struct {
	PullRequests []pullRequest `json:"nodes"`
	Total        int           `json:"totalCount"`
}

type assignees struct {
	TotalCount int `json:"totalCount"`
}
