// Command build-countries regenerates data/countries.json, the country dataset
// embedded in the API. It merges two open datasets (both ODbL-1.0) by ISO 3166-1
// alpha-2 code:
//
//   - mledoze/countries: names, alternative spellings, capital, area,
//     languages, currencies
//   - dr5hn/countries-states-cities-database: population and time zones
//
// Run from the repository root:
//
//	go run ./scripts/build-countries
//
// Sources are pinned to commits so the output is reproducible; bump the
// commit constants to refresh the data.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	mledozeCommit = "fa4f62182d19ff4f8b325cde6387c9cd9d80955a"
	dr5hnCommit   = "20de077b176b00ba0fc0f298a55caa1dc6337a86"

	mledozeURL = "https://raw.githubusercontent.com/mledoze/countries/" + mledozeCommit + "/countries.json"
	dr5hnURL   = "https://raw.githubusercontent.com/dr5hn/countries-states-cities-database/" + dr5hnCommit + "/json/countries.json"

	outputPath = "data/countries.json"
)

type mledozeCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	AltSpellings []string                     `json:"altSpellings"`
	Capital      []string                     `json:"capital"`
	Area         float64                      `json:"area"`
	Languages    map[string]string            `json:"languages"`
	Currencies   map[string]map[string]string `json:"currencies"`
	CCA2         string                       `json:"cca2"`
}

type dr5hnCountry struct {
	ISO2       string `json:"iso2"`
	Population int64  `json:"population"`
	Timezones  []struct {
		GMTOffsetName string `json:"gmtOffsetName"`
	} `json:"timezones"`
}

// Country is one record in data/countries.json.
type Country struct {
	Name         string                       `json:"name"`
	OfficialName string                       `json:"official_name"`
	AltSpellings []string                     `json:"alt_spellings"`
	CCA2         string                       `json:"cca2"`
	Capital      []string                     `json:"capital"`
	Population   int64                        `json:"population"`
	Area         float64                      `json:"area"`
	Languages    map[string]string            `json:"languages"`
	Timezones    []string                     `json:"timezones"`
	Currencies   map[string]map[string]string `json:"currencies"`
}

func fetch(url string, v any) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func main() {
	var base []mledozeCountry
	if err := fetch(mledozeURL, &base); err != nil {
		log.Fatalf("mledoze/countries: %v", err)
	}
	var extra []dr5hnCountry
	if err := fetch(dr5hnURL, &extra); err != nil {
		log.Fatalf("dr5hn/countries-states-cities-database: %v", err)
	}

	byISO2 := make(map[string]dr5hnCountry, len(extra))
	for _, c := range extra {
		byISO2[strings.ToUpper(c.ISO2)] = c
	}

	countries := make([]Country, 0, len(base))
	missing := 0
	for _, b := range base {
		c := Country{
			Name:         b.Name.Common,
			OfficialName: b.Name.Official,
			AltSpellings: b.AltSpellings,
			CCA2:         strings.ToUpper(b.CCA2),
			Capital:      b.Capital,
			Area:         b.Area,
			Languages:    b.Languages,
			Currencies:   b.Currencies,
			Timezones:    []string{},
		}
		if e, ok := byISO2[c.CCA2]; ok {
			c.Population = e.Population
			for _, tz := range e.Timezones {
				if tz.GMTOffsetName != "" && !contains(c.Timezones, tz.GMTOffsetName) {
					c.Timezones = append(c.Timezones, tz.GMTOffsetName)
				}
			}
		} else {
			missing++
		}
		countries = append(countries, c)
	}

	sort.Slice(countries, func(i, j int) bool { return countries[i].Name < countries[j].Name })

	out, err := json.MarshalIndent(countries, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outputPath, append(out, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d countries to %s (%d without population/time-zone data)\n", len(countries), outputPath, missing)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
