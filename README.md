## Overview

`gorate` provides an atomic, thread-safe rate limiting solution suitable for distributed systems. It leverages Redis Lua scripts to ensure that rate-limiting checks and token refills are performed in a single atomic operation, preventing race conditions in high-concurrency environments.

### Key Features

- **Token Bucket Algorithm:** Supports flexible burst capacity and steady-state refill rates.
- **Atomic Operations:** Uses Redis Lua scripting for consistency across multiple service instances.
- **Middleware Support:** Includes Go HTTP middleware for easy integration with existing web servers.
- **Extensible Storage:** Built around a `Store` interface, allowing for alternative backends (e.g., in-memory, Postgres) if needed.

## Architecture

The project is organized into several internal packages:

- **`internal/bucket`**: Core logic for the token bucket.
- **`internal/store`**: Storage abstractions and the Redis implementation.
- **`internal/middleware`**: HTTP middleware for rate limiting by client IP.
- **`cmd/server`**: Example server implementation.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.21 or later.
- [Redis](https://redis.io/docs/getting-started/) instance running locally or accessible via network.

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/nabsk911/gorate.git
cd gorate
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Run the server

Ensure Redis is running on `localhost:6379`, then start the server:

```bash
go run cmd/server/main.go
```

The server will start listening on `:8080`.

## Configuration

The default configuration in `main.go` sets up a bucket with:

- **Capacity:** 5 tokens (allows a burst of 5 requests).
- **Refill Rate:** 1 token per second.

You can modify these parameters in `cmd/server/main.go`:

```go
b := bucket.NewBucket(rdb, 10, 2, time.Second)
```

## Usage

### Testing the Rate Limit

You can use `curl` to test the rate limiter:

```bash
# First 5 requests will succeed
for i in {1..6}; do curl -i http://localhost:8080; done
```

When the limit is reached, you will receive a `429 Too Many Requests` response:

```json
{ "error": "rate limit exceeded" }
```

## References

- [Redis Rate Limiter Pattern (Go Implementation Guide)](https://redis.io/docs/latest/develop/use-cases/rate-limiter/go/)
