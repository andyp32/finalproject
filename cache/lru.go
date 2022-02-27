package cache

import "fmt"

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
	// current int
}

type Node struct {
	next     *Node
	previous *Node
	key      string
	value    []byte
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	lru := new(LRU)
	lru.limit = limit
	// lru.current = 0
	// lru.storage = byte[limit]
	lru.location = make(map[string]*Node)
	// lru.queue = make([]int, limit)
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
	val, ok := lru.location[key]
	if ok {
		if val.next != nil && val.previous != nil {
			val.next.previous = val.previous
			val.previous.next = val.next
			val.next = lru.root
			lru.root.previous = val
			lru.root = val
			lru.hits++
		} else if val.next != nil {
			lru.root = val
			lru.hits++
		} else if val.previous != nil {
			val.previous.next = nil
			val.next = lru.root
			lru.root.previous = val
			lru.root = val
			lru.hits++
		}

	} else {
		lru.misses++
	}

	// get fifo.storge[val] from storage

	return val.value, ok
}

// Pop least recently used value
func (lru *LRU) Pop() string {
	current := lru.lru_node
	value := lru.lru_node.key
	lru.lru_node = lru.lru_node.previous
	lru.lru_node.next = nil

	current.key = ""
	current.next = nil
	current.previous = nil
	current = nil
	return value
}

// Pop a given key
func (lru *LRU) PopKey(key string) string {
	value := ""
	current := lru.root
	for current != nil {
		current = current.next
		if key == current.key && current == lru.lru_node {
			lru.Pop()
		} else if key == current.key {
			value = current.key
			current.previous.next = current.next
			current.next.previous = current.previous

			current.next = lru.root
			current.previous = nil
		}
	}

	return value
}

// Add a given key to beginning of lru
func (lru *LRU) AddKey(key string, value []byte) *Node {
	current := new(Node)
	current.key = key
	current.value = value
	current.next = lru.root
	lru.root.previous = current
	lru.root = current
	return current
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	if val, ok := lru.location[key]; ok {
		lru.inUse += -len(key) - len(lru.location[key].value)

		delete(lru.location, key)
		lru.PopKey(key)
		lru.hits++
		lru.numBindings--
		return val.value, ok
	} else {
		lru.misses++
		return nil, false
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	size := len(key) + len(value)
	if size > lru.limit {
		return false
	}
	if lru.RemainingStorage() >= size {
		new_node := lru.AddKey(key, value)
		fmt.Println("Added")
		lru.location[key] = new_node
		fmt.Println("Mapped")

		lru.numBindings++
		lru.inUse += size
		return true
	} else {
		for lru.RemainingStorage() < size {
			lru_value := lru.Pop()
			lru.inUse += -len(lru_value) - len(lru.location[lru_value].value)
			delete(lru.location, lru_value)
		}
		new_node := lru.AddKey(key, value)
		lru.location[key] = new_node
		lru.inUse += size
		return true

	}
	// return false
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
