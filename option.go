package shuttle

import "time"

// Option you know
type Option func(pool *dataPool)

// WithFunc register shuttle function
func WithFunc(f ShuttleFunc) Option {
	return func(pool *dataPool) {
		pool.shuttleFunc = f
	}
}

// WithDuration set duration
func WithDuration(duration time.Duration) Option {
	return func(pool *dataPool) {
		pool.duration = duration
	}
}

// WithMaxSize set the maxSize
func WithMaxSize(size int) Option {
	return func(pool *dataPool) {
		pool.maxSize = size
	}
}
