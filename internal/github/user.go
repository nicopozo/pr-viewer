package github

type user struct {
	Username string `json:"login"`
}

type viewerResponse struct {
	Data struct {
		Viewer user `json:"viewer"`
	} `json:"data"`
}
