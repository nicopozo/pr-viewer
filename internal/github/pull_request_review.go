package github

const (
	pullRequestReviewStateCommented        = "COMMENTED"
	pullRequestReviewStateChangesRequested = "CHANGES_REQUESTED"
	pullRequestReviewStateApproved         = "APPROVED"
)

type PullRequestReview struct {
	Author User   `json:"author"`
	State  string `json:"state"`
}

// PullRequestReviewList represents a PullRequestReviewConnection
type PullRequestReviewList struct {
	PullRequestReview []PullRequestReview `json:"nodes"`
}
