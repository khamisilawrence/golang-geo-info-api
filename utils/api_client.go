package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"geo-info-api/data"
)

// ErrCountryNotFound is returned when no country matches the query.
var ErrCountryNotFound = errors.New("country not found")

// Country is one record of the embedded dataset (data/countries.json).
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

var (
	countriesOnce sync.Once
	countries     []Country
	countriesErr  error
)

// Countries returns the embedded dataset, parsing it on first use.
func Countries() ([]Country, error) {
	countriesOnce.Do(func() {
		countriesErr = json.Unmarshal(data.CountriesJSON, &countries)
		if countriesErr == nil && len(countries) == 0 {
			countriesErr = errors.New("embedded country dataset is empty")
		}
	})
	return countries, countriesErr
}

// FindCountry looks a country up by name. It tries, in order and ignoring
// case: the common name, the official name, alternative spellings (including
// ISO codes), then the first common name that starts with or contains the query.
func FindCountry(query string) (*Country, error) {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil, ErrCountryNotFound
	}

	list, err := Countries()
	if err != nil {
		return nil, err
	}

	matchers := []func(c *Country) bool{
		func(c *Country) bool { return strings.ToLower(c.Name) == q },
		func(c *Country) bool { return strings.ToLower(c.OfficialName) == q },
		func(c *Country) bool {
			for _, alt := range c.AltSpellings {
				if strings.ToLower(alt) == q {
					return true
				}
			}
			return false
		},
		func(c *Country) bool { return strings.HasPrefix(strings.ToLower(c.Name), q) },
		func(c *Country) bool { return strings.Contains(strings.ToLower(c.Name), q) },
	}

	for _, match := range matchers {
		for i := range list {
			if match(&list[i]) {
				return &list[i], nil
			}
		}
	}
	return nil, ErrCountryNotFound
}

// FetchCountryData returns the country matching the query, using the cache
// when possible. The result keeps the shape the handler used with the old
// upstream API, so the API response is unchanged.
func FetchCountryData(country string) (map[string]interface{}, error) {
	key := strings.ToLower(strings.TrimSpace(country))

	if Cache != nil {
		if cachedData, found := Cache.Get(key); found {
			return cachedData.(map[string]interface{}), nil
		}
	}

	c, err := FindCountry(key)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"name":       map[string]interface{}{"common": c.Name, "official": c.OfficialName},
		"capital":    c.Capital,
		"population": c.Population,
		"area":       c.Area,
		"languages":  c.Languages,
		"timezones":  c.Timezones,
		"currencies": c.Currencies,
		"flags":      map[string]interface{}{"png": flagURL(c.CCA2)},
	}

	if Cache != nil {
		Cache.Set(key, result, 1)
	}
	return result, nil
}

func flagURL(cca2 string) string {
	if cca2 == "" {
		return ""
	}
	return fmt.Sprintf("https://flagcdn.com/w320/%s.png", strings.ToLower(cca2))
}
