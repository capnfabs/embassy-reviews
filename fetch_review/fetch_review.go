/*
Suggested usage:

go install ./fetch_review && cat data/sample_fetch_embassies_response.txt | \
PLACES_API_KEY=`cat places_api_key.txt` fetch_review
*/

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flagAll := flag.Bool("all", false, "Include all reviews for all input places")
	flag.Parse()
	apiKey := os.Getenv(`PLACES_API_KEY`)
	if apiKey == "" {
		panic("Requires a valid places API key")
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Err() != nil {
			panic(sc.Err())
		}
		place := sc.Text()
		reviewResponse, err := fetchReviewsForPlace(apiKey, place)
		if err != nil {
			panic(err)
		}
		if reviewResponse.Status != "OK" {
			panic(fmt.Errorf("Bad response from Places API: %s", reviewResponse.Status))
		}
		for _, r := range shuffle(filter(reviewResponse.Result.Reviews)) {
			txt := fmt.Sprintf(
				"%s %s",
				strings.Repeat("★", r.Rating),
				strings.TrimSpace(r.Text))
			fmt.Println(txt)
			if !*flagAll {
				return
			}
		}
	}
}

// Skips non-english reviews, and reviews without text.
func filter(rs []review) []review {
	var out []review
	for _, r := range rs {
		if r.Language == "en" && r.Text != "" {
			out = append(out, r)
		}
	}
	return out
}

func shuffle(rs []review) []review {
	dest := make([]review, len(rs))
	perm := rand.Perm(len(rs))
	for i, v := range perm {
		dest[v] = rs[i]
	}
	return dest
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
	//log.Println(url.String())
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
