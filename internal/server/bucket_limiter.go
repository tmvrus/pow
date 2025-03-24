package server

import "sync"

type bucketLimiter struct {
	bucket map[string]int
	limit  int
	lock   sync.Mutex
}

func newBucketLimiter(size, limit int) *bucketLimiter {
	return &bucketLimiter{
		bucket: make(map[string]int, size),
		limit:  limit,
	}
}

func (l *bucketLimiter) acquire(key string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	current := l.bucket[key]
	if current < l.limit {
		l.bucket[key]++
		return true
	}
	return false
}

func (l *bucketLimiter) release(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.bucket[key]--
}
