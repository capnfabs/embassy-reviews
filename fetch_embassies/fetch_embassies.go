// Suggested Usage:
// go install ./fetch_embassies && PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Given LatLngs, fetches Places API place_ids for embassies in 50km radius around those LatLngs.
func main() {
	apiKey := os.Getenv(`PLACES_API_KEY`)
	if apiKey == "" {
		panic("Requires a valid places API key")
	}
	var embassies []place
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Err() != nil {
			panic(sc.Err())
		}
		vals := strings.Split(sc.Text(), ",")
		if len(vals) != 2 {
			panic("Expected a pair of coords")
		}
		lat := mustFloat32(vals[0])
		lng := mustFloat32(vals[1])
		page := ""
		count := 0
		for {
			// Do a cheeky sleep to prevent pages not being ready yet.
			if page != "" {
				time.Sleep(3 * time.Second)
			}
			resp, err := fetchEmbassiesNearLatLng(apiKey, page, lat, lng)
			if err != nil {
				log.Println(err)
			}
			if !(resp.Status == "OK" || resp.Status == "ZERO_RESULTS") {
				log.Println(resp.Status)
			}
			count += len(resp.Results)
			for _, r := range resp.Results {
				fmt.Println(r.PlaceID)
			}
			// TODO: convert to map on place ID.
			embassies = append(embassies, resp.Results...)
			page = resp.PageToken
			if page == "" {
				break
			}
		}
		log.Printf("Obtained %d results from %f,%f\n", count, lat, lng)
	}
}

func mustFloat32(str string) float32 {
	x, err := strconv.ParseFloat(strings.TrimSpace(str), 32)
	if err != nil {
		panic(err)
	}
	return float32(x)
}

func fetchEmbassiesNearLatLng(apiKey string, pageToken string, lat, lng float32) (placesResponse, error) {
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
	//log.Println(url.String())
	res, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var resp placesResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&resp)
	return resp, err
}
