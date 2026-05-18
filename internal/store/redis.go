package store

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// luaScript implements the token bucket algorithm atomically in Redis.
// It handles token refills based on time passed since the last refill.
const luaScript = `
local key             = KEYS[1]
local capacity        = tonumber(ARGV[1])
local refill_rate     = tonumber(ARGV[2])
local refill_interval = tonumber(ARGV[3])
local now             = tonumber(ARGV[4])
 
-- Get current bucket state: tokens and last refill timestamp.
local bucket      = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens      = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])
 
-- If bucket doesn't exist, initialize it with full capacity.
if tokens == nil then
    tokens      = capacity
    last_refill = now
end
 
-- Calculate how many tokens should be added based on time elapsed.
local time_passed = now - last_refill
local refills     = math.floor(time_passed / refill_interval)
 
-- If at least one refill interval has passed, update tokens.
if refills > 0 then
    tokens      = math.min(capacity, tokens + (refills * refill_rate))
    last_refill = last_refill + (refills * refill_interval)
end
 
-- Attempt to consume 1 token.
local allowed = 0
if tokens >= 1 then
    tokens  = tokens - 1
    allowed = 1
end
 
-- Persist the updated state back to Redis.
redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
 
return {allowed, tokens}
`

type RedisStore struct {
	client *redis.Client
	once   sync.Once
	sha    string
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisStore) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Take executes the Lua script in Redis to consume a token.
func (r *RedisStore) Take(ctx context.Context, ip string, capacity int, refillRate int, refillInterval time.Duration) (Result, error) {
	// sync.Once ensures the Lua script is loaded into Redis only once.
	// Redis returns a SHA1 hash which is cached for future EvalSha calls.
	r.once.Do(func() {
		r.sha, _ = r.client.ScriptLoad(ctx, luaScript).Result()
	})

	// EvalSha executes the cached Lua script using its SHA hash
	// instead of sending the entire script every request.
	res, err := r.client.EvalSha(ctx, r.sha, []string{ip}, capacity, refillRate, refillInterval.Milliseconds(), time.Now().UnixMilli()).Int64Slice()
	if err != nil {
		return Result{Allowed: false}, err
	}

	// res[0] is the allowed flag (1 or 0), res[1] is the remaining tokens.
	return Result{Allowed: res[0] == 1, Remaining: float64(res[1])}, nil
}
