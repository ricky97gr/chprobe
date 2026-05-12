package middleware

import (
	"sync"
	"time"
)

type RateLimiter struct {
	records map[string]int64
	mu      sync.Mutex
}

var GlobalDownloadRateLimiter = &RateLimiter{
	records: make(map[string]int64),
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().Unix()
	lastTime, exists := rl.records[key]
	if !exists || now-lastTime >= 60 {
		rl.records[key] = now
		return true
	}
	return false
}
