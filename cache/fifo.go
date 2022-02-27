package cache

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	front       int
	back        int
	limit       int
	inUse       int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	location map[string][]byte
	// storage  *byte
	queue   []string
	hits    int
	misses  int
	current int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := new(FIFO)
	fifo.front = 0
	fifo.back = 0
	fifo.limit = limit
	fifo.current = 0
	// fifo.storage = byte[limit]
	fifo.location = make(map[string][]byte)
	fifo.queue = make([]string, limit)
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
	return fifo.limit - fifo.inUse
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

// Pop first in first out value
func (fifo *FIFO) Pop() string {
	first := fifo.queue[0]
	fifo.PopKey(first)
	return first
}

// Pop first in first out value
func (fifo *FIFO) PopKey(key string) string {

	write := false
	value := ""
	for i, v := range fifo.queue {
		if key == v {
			value = v
			write = true
		} else if write {
			fifo.queue[i-1] = v
		}
	}

	return value
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	if val, ok := fifo.location[key]; ok {
		delete(fifo.location, key)
		fifo.PopKey(key)
		// get fifo.storge[val] from storage
		fifo.hits++
		fifo.numBindings--
		return val, ok
	} else {
		fifo.misses++
		return nil, false
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {

	if fifo.RemainingStorage() > 0 {
		fifo.location[key] = value
		fifo.numBindings++
		fifo.queue[fifo.numBindings] = key
		return true
	} else {
		first := fifo.Pop()
		delete(fifo.location, first)

		fifo.location[key] = value
		fifo.queue[fifo.numBindings] = key

	}

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
