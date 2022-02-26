package cache

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	front    int
	back     int
	limit    int
	inUse     int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	location map[string]{int, int}
	storage  *byte
	hits     int
	misses   int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := new(FIFO)
	fifo.front = 0
	fifo.back = 0
	fifo.limit = limit
	fifo.curr = 0
	fifo.storage = byte[limit]
	fifo.location = make(map[string]int)
	fifo.hits = 0
	fifo.misses = 0
	return fifo
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.limit
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.Limit - fifo.inUse
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	val, ok := fifo.location[key]
	if ok {
		fifo.hits++
	} else {
		fifo.misses++
	}

	// get fifo.storge[val] from storage

	return val, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	if val, ok := fifo.location[key]; ok {
		delete(fifo.location, key)

		// get fifo.storge[val] from storage
		fifo.hits++
		return val, ok
	} else {
		fifo.misses++
		return nil, false
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {
	return false
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return fifo.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return &Stats{Hits: fifo.hits, Misses: fifo.misses}
}
