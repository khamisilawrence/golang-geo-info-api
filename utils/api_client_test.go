package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountriesDatasetLoads(t *testing.T) {
	list, err := Countries()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 240, "embedded dataset looks truncated")
}

func TestFetchCountryData(t *testing.T) {
	InitCache() // Ensure cache is initialized
	country := "Germany"

	// First call: look up in the embedded dataset
	data, err := FetchCountryData(country)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "Germany", data["name"].(map[string]interface{})["common"])

	// Second call: may come from cache; must be equal either way
	Cache.Wait()
	cachedData, err := FetchCountryData(country)
	assert.NoError(t, err)
	assert.Equal(t, data, cachedData)
}

func TestFindCountryMatching(t *testing.T) {
	cases := map[string]string{
		"jamaica":                     "Jamaica",
		"  JAMAICA ":                  "Jamaica",
		"Federal Republic of Germany": "Germany",
		"JM":                          "Jamaica",
		"united king":                 "United Kingdom",
	}
	for query, want := range cases {
		c, err := FindCountry(query)
		if assert.NoError(t, err, query) {
			assert.Equal(t, want, c.Name, query)
		}
	}
}

func TestFindCountryNotFound(t *testing.T) {
	_, err := FindCountry("Atlantis")
	assert.ErrorIs(t, err, ErrCountryNotFound)

	_, err = FindCountry("")
	assert.ErrorIs(t, err, ErrCountryNotFound)
}
