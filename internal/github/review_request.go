package github

type reviewRequest struct {
	RequestedReviewer user `json:"requestedReviewer"`
}

// reviewRequestList represents a ReviewRequestConnection
type reviewRequestList struct {
	ReviewRequests []reviewRequest `json:"nodes"`
}
