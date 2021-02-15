package github

type ReviewRequest struct {
	RequestedReviewer User `json:"requestedReviewer"`
}

// ReviewRequestList represents a ReviewRequestConnection
type ReviewRequestList struct {
	ReviewRequests []ReviewRequest `json:"nodes"`
}
