/*
Suggested usage:

go install ./fetch_review && cat data/sample_fetch_embassies_response.txt | \
PLACES_API_KEY=`cat places_api_key.txt` fetch_review
*/

package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/capnfabs/embassyreviews/reviews"
)

func main() {
	var placesList = make([]reviews.PlaceDetails, 0)
	apiKey := os.Getenv(`PLACES_API_KEY`)
	if apiKey == "" {
		panic("Requires a valid places API key")
	}
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	signal.Notify(sigInt, syscall.SIGTERM)

	sc := bufio.NewScanner(os.Stdin)
	lastPlace := ""
	errCount := 0
input:
	for sc.Scan() {
		if sc.Err() != nil {
			log.Print(sc.Err())
			break input
		}
		// Check that we haven't been interrupted
		select {
		case <-sigInt:
			log.Println("Got SIGINT, shutting down")
			break input
		default:
		}

		place := sc.Text()
		reviewResponse, err := fetchReviewsForPlace(apiKey, place)

		if err != nil {
			log.Print(err)
			// Hope that network errors go away after 5 seconds
			time.Sleep(5 * time.Second)
			// 10 network failures in a row before giving up.
			errCount++
			if errCount > 10 {
				break input
			}
			continue
		}
		errCount = 0

		if reviewResponse.Status != "OK" {
			log.Printf("Bad response from Places API: %s", reviewResponse.Status)
			break input
		}

		lastPlace = place
		log.Printf("Fetched details for %s", lastPlace)
		placesList = append(placesList, reviewResponse.Result)
	}
	log.Print("Finishing! Made it all the way to ", lastPlace)
	if err := outputToJSON(placesList); err != nil {
		panic(err)
	}
}

func outputToJSON(val interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(val)
}

func fetchReviewsForPlace(apiKey, placeID string) (reviews.ReviewResponse, error) {
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
	var resp reviews.ReviewResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&resp)
	return resp, err
}
