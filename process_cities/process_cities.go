package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var geonameFields = []string{
	`geonameid`,
	`name`,
	`asciiname`,
	`alternatenames`,
	`latitude`,
	`longitude`,
	`feature_class`,
	`feature_code`,
	`country_code`,
	`cc2`,
	`admin1_code`,
	`admin2_code`,
	`admin3_code`,
	`admin4_code`,
	`population`,
	`elevation`,
	`dem`,
	`timezone`,
	`modification_date`,
}

var countryInfoFields = []string{
	`iso`,
	`iso3`,
	`iso-numeric`,
	`fips`,
	`country`,
	`capital`,
	`area`,
	`population`,
	`continent`,
	`tld`,
	`currencycode`,
	`currencyname`,
	`phone`,
	`postal_code_format`,
	`postal_code_regex`,
	`languages`,
	`geonameid`,
	`neighbours`,
	`equivalentfipscode`,
}

type latLng struct {
	lat float32
	lng float32
}

type city struct {
	country string
	city    string
}

// Double quotes are used extensively in the set, including at the start of fields, but they're
// never used to wrap a field. We could roll our own, more general, CSV parser, or we could
// just replace characters on the fly. We're going to go with a custom parser because everything's
// pretty well formatted.
func parseGeonamesFile(ioreader io.Reader) [][]string {
	s := bufio.NewScanner(ioreader)
	var records [][]string
	count := 0
	for lineno := 0; s.Scan(); lineno++ {
		l := s.Text()
		// Skip blank lines
		if l == "" || l[0] == '#' {
			//log.Println("Skipping line ", lineno)
			continue
		}
		fields := strings.Split(l, "\t")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		if count == 0 {
			count = len(fields)
		} else if count != len(fields) {
			panic(fmt.Sprintf("Line %d: Got %d fields, expected %d\n", lineno, len(fields), count))
		}
		records = append(records, fields)
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return records
}

func loadCityLatLngs(ioreader io.Reader) map[city]latLng {
	records := parseGeonamesFile(ioreader)
	x := make(map[city]latLng)
	for _, rec := range records {
		f := mapSlices(geonameFields, rec)
		city := city{
			country: f[`country_code`],
			city:    f[`asciiname`],
		}
		x[city] = latLng{
			lat: mustFloat32(f[`latitude`]),
			lng: mustFloat32(f[`longitude`]),
		}
	}
	return x
}

func loadCapitals(ioreader io.Reader) []city {
	records := parseGeonamesFile(ioreader)
	var x []city
	for _, rec := range records {
		f := mapSlices(countryInfoFields, rec)
		cap := f[`capital`]
		if cap == "" {
			continue
		}
		city := city{
			country: f[`iso`],
			city:    f[`capital`],
		}
		x = append(x, city)
	}
	return x
}

func mustFloat32(str string) float32 {
	x, err := strconv.ParseFloat(str, 32)
	if err != nil {
		panic(err)
	}
	return float32(x)
}

func mapSlices(names, values []string) map[string]string {
	if len(names) != len(values) {
		panic(fmt.Sprintf("Expected both inputs to have same length, got len(names)=%d len(values)=%d", len(names), len(values)))
	}
	x := make(map[string]string)
	for i, name := range names {
		x[name] = values[i]
	}
	return x
}

func main() {
	r, err := os.Open("data/cities15000.txt")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	cityLatLngs := loadCityLatLngs(r)
	log.Printf("Loaded %d cities w/ coordinates\n", len(cityLatLngs))

	r, err = os.Open("data/countryInfo.txt")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	capitals := loadCapitals(r)
	log.Printf("Loaded %d capital cities", len(capitals))
	found := 0
	missed := 0
	for _, c := range capitals {
		l, ok := cityLatLngs[c]
		if !ok {
			log.Printf("Couldn't find %v\n", c)
			missed++
			continue
		}
		log.Printf("%s %s: %f,%f\n", c.country, c.city, l.lat, l.lng)
		found++
	}
	log.Printf("Found: %d Missed %d\n", found, missed)
}
