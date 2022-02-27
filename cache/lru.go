package cache

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	limit       int
	inUse       int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	location map[string]*Node
	// storage  *byte
	hits     int
	misses   int
	root     *Node
	lru_node *Node
}

// type Node struct {
// 	next     *Node
// 	previous *Node
// 	key      string
// 	value    []byte
// }

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	lru := new(LRU)
	lru.limit = limit
	lru.location = make(map[string]*Node)
	lru.root = new(Node)
	lru.lru_node = lru.root
	lru.hits = 0
	lru.misses = 0
	return lru
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.limit
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.limit - lru.inUse

}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	return false
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
