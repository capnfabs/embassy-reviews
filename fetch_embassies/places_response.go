package main

type placesResponse struct {
	PageToken string  `json:"next_page_token"`
	Results   []place `json:"results"`
	Status    string  `json:"status"`
}

type place struct {
	PlaceID string `json:"place_id"`
	// We don't need this, but we may as well not throw away the data.
	Name     string   `json:"name"`
	Geometry geometry `json:"geometry"`
	Rating   float32  `json:"rating,omitempty"`
}

type geometry struct {
	Location latLng `json:"location"`
}

type latLng struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}
