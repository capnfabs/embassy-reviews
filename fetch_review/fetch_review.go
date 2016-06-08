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
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

func main() {
	var reviewList []string
	rand.Seed(time.Now().UnixNano())
	flagAll := flag.Bool("all", false, "Include all reviews for all input places")
	flag.Parse()
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
			log.Println("Got sigint, shutting down")
			break input
		default:
		}

		place := sc.Text()
		reviewResponse, err := fetchReviewsForPlace(apiKey, place)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second)
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
		for _, r := range shuffle(filter(reviewResponse.Result.Reviews)) {
			// Choose 140 - 30 (for reviews *, links, hyphens, spaces.)
			reviewText := limitChooseSentence(r.Text, 110)
			txt := fmt.Sprintf(
				"%s %s - %s",
				strings.Repeat("★", r.Rating),
				reviewText,
				reviewResponse.Result.URL)
			log.Println(txt)
			reviewList = append(reviewList, txt)
			if !*flagAll {
				break input
			}
		}
	}
	log.Print("Finishing! Made it all the way to ", lastPlace)
	if len(reviewList) > 0 {
		shuffleStrings(reviewList)
		if err := outputToJSON(reviewList); err != nil {
			panic(err)
		}
	}
}

func outputToJSON(val interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(val)
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

func shuffleStrings(in []string) []string {
	dest := make([]string, len(in))
	perm := rand.Perm(len(in))
	for i, v := range perm {
		dest[v] = in[i]
	}
	return dest
}

func limitChooseSentence(text string, maxlength int) string {
	if len(text) < maxlength {
		return text
	}
	sentences := removeEmpty(strings.Split(text, "."))
	order := rand.Perm(len(sentences))
	var chosen []int
	sum := 0
	for _, idx := range order {
		// +1 for the full stop + space.
		lenItem := len(sentences[idx]) + 2
		if sum+lenItem <= maxlength {
			chosen = append(chosen, idx)
			sum += lenItem
		}
	}
	// Sort the sentences again to put them in the right order
	sort.Ints(chosen)
	retVal := ""
	// chosen so that it won't be triggered first time
	last := -5
	for _, idx := range chosen {
		if last == -5 {
			// first time, just append as is
			retVal += sentences[idx]
		} else if last == idx-1 {
			// Strings are right next to each other! Don't try to put a …
			retVal += ". " + sentences[idx]
		} else {
			// There's a gap! Praise the Codepoints in heaven above for the fact that an ellipsis is one
			// character.
			retVal += "… " + sentences[idx]
		}
		last = idx
	}
	return retVal
}

func removeEmpty(in []string) []string {
	var out []string
	for _, val := range in {
		if val != "" {
			out = append(out, strings.TrimSpace(val))
		}
	}
	return out
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
