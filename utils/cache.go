package utils

import (
	"log"

	"github.com/dgraph-io/ristretto"
)

var Cache *ristretto.Cache

func InitCache() {
	var err error
	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // Number of keys to track frequency
		MaxCost:     1 << 30, // Maximum memory usage (1GB)
		BufferItems: 64,      // Number of keys per buffer
	})
	if err != nil {
		log.Fatalf("Error initializing cache: %v", err)
	}
}
