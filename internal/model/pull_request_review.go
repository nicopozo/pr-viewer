package model

const (
	pullRequestReviewStateCommented        = "COMMENTED"
	pullRequestReviewStateChangesRequested = "CHANGES_REQUESTED"
	pullRequestReviewStateApproved         = "APPROVED"
)

type PullRequestReview struct {
	Author string `json:"author"`
	State  string `json:"state"`
}
