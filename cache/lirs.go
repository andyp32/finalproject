package cache

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type LIRS struct {
	limit       int
	inUse       int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	// location map[string]*Node
	S      *FIFO
	Q      *FIFO
	hits   int
	misses int
	// current int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewLIRS(limit int) *LIRS {
	lirs := new(LIRS)
	lirs.S = NewFifo(limit)
	lirs.Q = NewFifo(limit)
	lirs.limit = limit
	lirs.inUse = 0
	lirs.numBindings = 0
	return lirs
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (lirs *LIRS) MaxStorage() int {
	return lirs.limit
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (lirs *LIRS) Get(key string) (value []byte, ok bool) {
	return
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lirs *LIRS) Set(key string, value []byte) bool {
	return true
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lirs *LIRS) Remove(key string) (value []byte, ok bool) {
	return
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (lirs *LIRS) RemainingStorage() int {
	return lirs.limit - lirs.inUse
}

// Len returns the number of bindings in the FIFO.
func (lirs *LIRS) Len() int {
	return lirs.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lirs *LIRS) Stats() *Stats {
	return &Stats{Hits: lirs.hits, Misses: lirs.misses}
}
