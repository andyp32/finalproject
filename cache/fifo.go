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
	queue  []string
	hits   int
	misses int
	// current int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := new(FIFO)
	fifo.front = 0
	fifo.back = 0
	fifo.limit = limit
	// fifo.current = 0
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
	first := fifo.queue[fifo.front]
	fifo.front = (fifo.front + 1) % fifo.limit
	return first

}

// Pop a given key
func (fifo *FIFO) PopKey(key string) string {

	write := false
	value := ""
	tot_range := fifo.limit + fifo.front
	length := fifo.limit
	for i := fifo.front; i < tot_range; i++ {

		if key == fifo.queue[i%length] {
			value = fifo.queue[i%length]
			write = true
		} else if write {
			fifo.queue[(i-1)%length] = fifo.queue[i%length]
		}
	}

	return value
}

// Add a given key to end of queue
func (fifo *FIFO) AddKey(key string) string {

	fifo.queue[fifo.back] = key
	fifo.back = (fifo.back + 1) % fifo.limit
	return key
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	if val, ok := fifo.location[key]; ok {
		fifo.inUse += -len(key) - len(fifo.location[key])

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
	size := len(key) + len(value)
	if size > fifo.limit {
		return false
	}
	if fifo.RemainingStorage() >= size {
		fifo.location[key] = value
		fifo.numBindings++
		fifo.AddKey(key)
		fifo.inUse += size
		// fifo.queue[fifo.numBindings] = key
		return true
	} else {
		for fifo.RemainingStorage() < size {
			first := fifo.Pop()
			fifo.inUse += -len(first) - len(fifo.location[first])
			delete(fifo.location, first)
		}
		fifo.location[key] = value
		fifo.AddKey(key)
		fifo.inUse += size
		return true

	}

	// return false
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return fifo.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return &Stats{Hits: fifo.hits, Misses: fifo.misses}
}
