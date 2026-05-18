package store

import (
	"context"
	"time"
)

type Result struct {
	Allowed   bool
	Remaining float64
}

type Store interface {
	// Take attempts to consume a token from the bucket for the given key.
	// It returns a Result indicating if the token was successfully taken.
	Take(ctx context.Context, key string, capacity int, refillRate int, refillInterval time.Duration) (Result, error)
}
