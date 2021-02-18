package model

import "time"

const (
	pullRequestReviewStateCommented        = "COMMENTED"
	pullRequestReviewStateChangesRequested = "CHANGES_REQUESTED"
	pullRequestReviewStateApproved         = "APPROVED"
)

type PullRequestReview struct {
	Author    string    `json:"author"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}
