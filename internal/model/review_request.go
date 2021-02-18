package model

const (
	ReviewRequestStatusPending          = "PENDING"
	ReviewRequestStatusChangesRequested = "CHANGES_REQUESTED"
	ReviewRequestStatusApproved         = "APPROVED"
	ReviewRequestStatusCommented        = "COMMENTED"
	ReviewRequestStatusDismissed        = "DISMISSED"
)

type ReviewRequest struct {
	RequestedReviewer string `json:"requested_reviewer"`
	State             string `json:"state"`
}
