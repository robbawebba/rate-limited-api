package rate

import (
	"time"
)

// Limiter represents an instance of a rate limiter
type Limiter struct {
	rate              int
	bucketSize        int
	remainingRequests int
	destroy           chan interface{}
}

// NewLimiter creates a new rate limiter with a bucket size of b (initially full).
// The bucket is refilled at a rate of r requests per second
func NewLimiter(r int, b int) *Limiter {
	l := &Limiter{
		r,
		b,
		r,
		make(chan interface{}),
	}
	go refreshLimiter(l)
	return l
}

// Allow determines if there are any "tokens" remaining in the token bucket
func (l *Limiter) Allow() bool {
	if l.remainingRequests > 0 {
		l.remainingRequests = l.remainingRequests - 1
		return true
	}
	return false
}

// StopLimiting stops a rate limiter from running
func (l *Limiter) StopLimiting() {
	l.destroy <- struct{}{}
}

func refreshLimiter(l *Limiter) {

	for {
		select {
		case <-l.destroy:
			return
		default:
		}
		diff := l.remainingRequests - l.bucketSize
		if diff > l.rate {
			l.remainingRequests = l.remainingRequests + l.rate
		} else {
			l.remainingRequests = l.remainingRequests + diff
		}
		// fmt.Printf("refreshed L: %+v", l)
		time.Sleep(1000 * time.Millisecond)
	}
}
