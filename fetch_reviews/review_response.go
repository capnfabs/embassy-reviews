package main

type reviewResponse struct {
	Result result `json:"result"`
	Status string `json:"status"`
}

type result struct {
	Reviews []review `json:"reviews"`
	// Just to check that it works
	Name string `json:"name"`
	URL  string `json:"url"`
}

type review struct {
	Rating   int    `json:"rating"`
	Text     string `json:"text"`
	Language string `json:"language"`
}
