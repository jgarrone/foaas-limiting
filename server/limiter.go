package server

import (
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Limiter interface {
	AllowRequestFrom(userid string) bool
}

type tokenBucketLimiter struct {
	cache *gocache.Cache
	mu    sync.Mutex
	limit int
}

// NewTokenBucketLimiter is an implementation of the Limiter interface that uses a key-value in-memory cache
// to keep track of user request counts.
// The cache automatically removes keys after a fixed interval of time, which is set to be the unit of time
// where requests should be limited.
func NewTokenBucketLimiter(limit int, window time.Duration) Limiter {
	return &tokenBucketLimiter{
		cache: gocache.New(window /* defaultExpiration */, window*2 /* cleanupInterval */),
		limit: limit,
	}
}

func (l *tokenBucketLimiter) AllowRequestFrom(userid string) bool {
	key := userid

	l.mu.Lock()
	defer l.mu.Unlock()

	tokens := l.limit

	val, ok := l.cache.Get(key)
	if ok {
		tokens = val.(int)
	}

	if tokens == 0 {
		return false
	}

	l.cache.SetDefault(key, tokens-1)

	return true
}
