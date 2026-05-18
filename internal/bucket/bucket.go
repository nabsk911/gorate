package bucket

import (
	"context"
	"time"

	"github.com/nabsk911/gorate/internal/store"
)

type Bucket struct {
	store          store.Store   // The underlying storage mechanism for bucket state.
	capacity       int           // Maximum number of tokens the bucket can hold.
	refillRate     int           // Number of tokens added to the bucket per refill interval.
	refillInterval time.Duration // Time interval between refills.
}

func NewBucket(s store.Store, capacity int, refillRate int, refillInterval time.Duration) *Bucket {
	return &Bucket{
		store:          s,
		capacity:       capacity,
		refillRate:     refillRate,
		refillInterval: refillInterval,
	}
}

// Allow checks if a request from the given IP address should be permitted.
// It interacts with the underlying store to decrement tokens and check availability.
func (b *Bucket) Allow(ctx context.Context, ip string) (store.Result, error) {
	return b.store.Take(ctx, ip, b.capacity, b.refillRate, b.refillInterval)
}
