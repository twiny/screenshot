package rate

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// client
type client struct {
	addr      string
	rate      *rate.Limiter
	createdAt time.Time
}

// Limiter
type Limiter struct {
	mu     *sync.RWMutex
	rate   rate.Limit
	bursts int
	store  map[string]*client
}

// NewRateLimiter
func NewLimiter(r float64, b int) *Limiter {
	return &Limiter{
		mu:     &sync.RWMutex{},
		rate:   rate.Limit(r),
		bursts: b,
		store:  map[string]*client{},
	}
}

// FindIPAddr
func (l *Limiter) FindIPAddr(ip string) *rate.Limiter {
	l.mu.RLock()
	defer l.mu.RUnlock()
	//
	c, found := l.store[ip]
	if !found {
		c = &client{
			addr:      ip,
			rate:      rate.NewLimiter(l.rate, l.bursts),
			createdAt: time.Now(),
		}
		l.store[ip] = c
	}

	return c.rate
}

// Clean
func (l *Limiter) Clean(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	//
	for k, v := range l.store {
		if time.Since(v.createdAt) > d {
			delete(l.store, k)
		}
	}
}

// Statistic
func (l *Limiter) Statistic() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()
	//
	stats := map[string]interface{}{
		"total": len(l.store),
	}
	return stats
}
