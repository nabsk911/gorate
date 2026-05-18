package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nabsk911/gorate/internal/bucket"
	"github.com/nabsk911/gorate/internal/middleware"
	"github.com/nabsk911/gorate/internal/store"
)

func main() {
	rdb := store.NewRedisStore("localhost:6379")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx); err != nil {
		log.Fatalf("redis unreachable: %v", err)
	}

	b := bucket.NewBucket(rdb, 5, 1, time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", middleware.RateLimit(b)(mux)); err != nil {
		log.Fatal(err)
	}
}
