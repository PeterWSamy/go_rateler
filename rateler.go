package gorateler

import (
	"math"
	"sync"
	"time"
)

type bucket struct {
	lastVisit time.Time
	tokens    float64
}

type safeRateler struct {
	mu       sync.Mutex
	rate     float64
	capacity int
	buckets  map[string]*bucket
}

func (rate *safeRateler) Allowed(ip string) bool {
	rate.mu.Lock()
	defer rate.mu.Unlock()

	now := time.Now()
	ipBucket, ok := rate.buckets[ip]

	if !ok {
		rate.buckets[ip] = &bucket{now, rate.rate}
		return true
	}
	timeElapsed := now.Sub(ipBucket.lastVisit).Seconds()
	ipBucket.tokens = math.Min(rate.rate, ipBucket.tokens+timeElapsed*rate.rate)
	if ipBucket.tokens < 1 {
		return false
	}
	ipBucket.tokens--
	return true
}

func NewRatelerTokenBucket(capacity int, rate float64) *safeRateler {
	return &safeRateler{
		rate:     rate,
		capacity: capacity,
		buckets:  make(map[string]*bucket),
	}
}
