package reviews

type ReviewResponse struct {
	Result PlaceDetails `json:"result"`
	Status string       `json:"status"`
}

type PlaceDetails struct {
	Reviews []Review `json:"reviews"`
	// Just to check that it works
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Review struct {
	Rating   int    `json:"rating"`
	Text     string `json:"text"`
	Language string `json:"language"`
}
