package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchCountryData calls the Restcountries API
func FetchCountryData(country string) (map[string]interface{}, error) {
	// Check cache first
	if cachedData, found := Cache.Get(country); found {
		return cachedData.(map[string]interface{}), nil
	}

	// Fetch from external API if not in cache
	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	// Return the first match
	if len(data) > 0 {
		Cache.Set(country, data[0], 1) // Cache the first result
		return data[0], nil
	}

	return nil, fmt.Errorf("no data found for country: %s", country)
}
