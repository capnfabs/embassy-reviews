package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/capnfabs/embassyreviews/reviews"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	dec := json.NewDecoder(os.Stdin)
	var placeDetails []reviews.PlaceDetails
	err := dec.Decode(&placeDetails)
	if err != nil {
		panic(err)
	}
	var formattedReviews = make([]string, 0)
	for _, place := range placeDetails {
		formattedReviews = append(formattedReviews, processPlace(place)...)
	}
	if len(formattedReviews) > 0 {
		// TODO: I don't think this works?
		formattedReviews = shuffleStrings(formattedReviews)
		if err := outputToJSON(formattedReviews); err != nil {
			panic(err)
		}
	}
}

func outputToJSON(val interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(val)
}

func processPlace(place reviews.PlaceDetails) []string {
	var output []string
	for _, r := range filter(place.Reviews) {
		// Choose 140 - 30 (for reviews *, links, hyphens, spaces.)
		reviewText := limitChooseSentence(r.Text, 110)
		if reviewText == "" {
			// Skip reviews that are empty once cleaned or shortened.
			continue
		}

		txt := fmt.Sprintf(
			"%s %s %s",
			strings.Repeat("★", r.Rating),
			reviewText,
			place.URL)
		log.Println(txt)
		output = append(output, txt)
	}
	return output
}

// Skips non-english reviews, and reviews without text.
func filter(rs []reviews.Review) []reviews.Review {
	var out []reviews.Review
	for _, r := range rs {
		if r.Language == "en" && r.Text != "" {
			out = append(out, r)
		}
	}
	return out
}

func shuffle(rs []reviews.Review) []reviews.Review {
	dest := make([]reviews.Review, len(rs))
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
		return strings.TrimSpace(text)
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
