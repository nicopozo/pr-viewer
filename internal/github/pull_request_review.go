package github

import "time"

const (
	pullRequestReviewStateCommented        = "COMMENTED"
	pullRequestReviewStateChangesRequested = "CHANGES_REQUESTED"
	pullRequestReviewStateApproved         = "APPROVED"
)

type pullRequestReview struct {
	Author    user      `json:"author"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PullRequestReviewList represents a PullRequestReviewConnection
type pullRequestReviewList struct {
	PullRequestReview []pullRequestReview `json:"nodes"`
}
