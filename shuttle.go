package shuttle

// New init a shuttle
func New(options ...Option) Pool {
	pool := newDataPool(options...)
	return pool
}
