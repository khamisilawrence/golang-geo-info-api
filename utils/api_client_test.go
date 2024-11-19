package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCountryData(t *testing.T) {
	InitCache() // Ensure cache is initialized
	country := "Germany"

	// First call: Fetch from API
	data, err := FetchCountryData(country)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "Germany", data["name"].(map[string]interface{})["common"])

	// Second call: Fetch from cache
	cachedData, err := FetchCountryData(country)
	assert.NoError(t, err)
	assert.Equal(t, data, cachedData)
}
