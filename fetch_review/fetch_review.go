/*
Suggested usage:

go install ./fetch_review && cat data/sample_fetch_embassies_response.txt | \
PLACES_API_KEY=`cat places_api_key.txt` fetch_review
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	apiKey := os.Getenv(`PLACES_API_KEY`)
	if apiKey == "" {
		panic("Requires a valid places API key")
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Err() != nil {
			panic(sc.Err())
		}
		reviewResponse, err := fetchReviewsForPlace(apiKey, sc.Text())
		if err != nil {
			log.Print(err)
		}
		fmt.Println(reviewResponse)
	}
}

func fetchReviewsForPlace(apiKey, placeID string) (reviewResponse, error) {
	url, err := url.Parse(`https://maps.googleapis.com/maps/api/place/details/json`)
	if err != nil {
		panic(err)
	}
	q := url.Query()
	q.Set(`key`, apiKey)
	q.Set(`placeid`, placeID)
	q.Set(`language`, `en`)
	url.RawQuery = q.Encode()
	log.Println(url.String())
	res, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var resp reviewResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&resp)
	return resp, err
}
