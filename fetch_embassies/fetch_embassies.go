// Suggested Usage:
// go install ./fetch_embassies && PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Given LatLngs, fetches Places API place_ids for embassies in 50km radius around those LatLngs.
func main() {

	apiKey := os.Getenv(`PLACES_API_KEY`)
	if apiKey == "" {
		panic("Requires a valid places API key")
	}
	fetchEmbassiesNearLatLng(apiKey, "", -35.283459, 149.128067)
}

func fetchEmbassiesNearLatLng(apiKey string, pageToken string, lat, lng float32) []string {
	url, err := url.Parse(`https://maps.googleapis.com/maps/api/place/nearbysearch/json`)
	if err != nil {
		panic(err)
	}
	q := url.Query()
	q.Set(`key`, apiKey)
	if pageToken != "" {
		q.Set(`pagetoken`, pageToken)
	} else {
		q.Set(`location`, fmt.Sprintf("%f,%f", lat, lng))
		q.Set(`radius`, `50000`)
		q.Set(`language`, `en`)
		q.Set(`type`, `embassy`)
	}
	url.RawQuery = q.Encode()
	fmt.Println(url.String())
	res, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(content))
	var ret []string
	return ret
}
