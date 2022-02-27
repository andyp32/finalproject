package cache

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	limit       int
	inUse       int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	location map[string]*Node
	front    *Node
	back     *Node
	hits     int
	misses   int
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	lru := new(LRU)
	lru.front = new(Node)
	lru.back = lru.front
	lru.limit = limit
	lru.location = make(map[string]*Node)
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
func (lru *LRU) PlaceNodeFront(current *Node) bool {
	if current == lru.front {
		// Do nothing
	} else if lru.numBindings == 2 {
		// current=back; current -> front; current = front; front -> next
		lru.back = current.previous

		current.next = lru.front
		lru.front = current
		lru.front.next = lru.back
		lru.back.previous = lru.front
		lru.front.previous = nil
		lru.back.next = nil

		// fmt.Println("======================================")
		// fmt.Println(lru.back)
		// fmt.Println("======================================")
		// fmt.Println(lru.front)
	} else {
		// current=middle; current -> front; current = front; front -> next
		current.next.previous = current.previous
		current.previous.next = current.next
		current.next = lru.front
		current.previous = nil
		lru.front = current
	}
	return true
}
func (lru *LRU) DeleteNode(current *Node) bool {

	if lru.numBindings == 1 {
	} else if current == lru.back {
		// current=back; current -> front; current = front; front -> next

		lru.back = current.previous
	} else if current == lru.front {
		lru.front = current.next
	} else {
		// current=middle; current -> front; current = front; front -> next
		current.next.previous = current.previous
		current.previous.next = current.next
	}

	lru.inUse += -len(current.key) - len(lru.location[current.key].value)
	delete(lru.location, current.key)
	current = nil
	lru.numBindings--

	return true
}
func (lru *LRU) CreateNode(key string, value []byte) bool {
	size := len(key) + len(value)

	if lru.numBindings == 0 {
		lru.front.key = key
		lru.front.value = value
		lru.location[key] = lru.front

	} else if lru.numBindings == 1 {
		new_node := new(Node)
		new_node.next = lru.front
		new_node.previous = nil
		new_node.value = value
		new_node.key = key
		new_node.size = size
		lru.front.previous = new_node
		lru.front = new_node
		lru.location[key] = new_node

		lru.back = lru.front.next
		lru.back.previous = lru.front
		lru.front.previous = nil
		lru.back.next = nil
	} else {
		new_node := new(Node)
		new_node.next = lru.front
		new_node.previous = nil
		new_node.value = value
		new_node.key = key
		new_node.size = size
		lru.front.previous = new_node
		lru.front = new_node

		lru.location[key] = new_node
	}

	lru.inUse += size
	lru.numBindings++
	return true

}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	val, ok := lru.location[key]
	if ok {
		lru.PlaceNodeFront(val)
		lru.hits++
		return val.value, ok

	} else {
		lru.misses++
		return nil, false
	}
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	if val, ok := lru.location[key]; ok {
		lru.DeleteNode(val)
		lru.hits++
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
	// if binding exists update old value, if size permits
	if val, ok := lru.location[key]; ok {
		bytes_diff := len(value) - len(val.value) // check space needed
		if lru.RemainingStorage() >= bytes_diff {
			lru.DeleteNode(val)
			lru.CreateNode(key, value)
			return true
		}
		return false

		// else if a existing binding is not found then need to create from scratch
	} else {

		if lru.RemainingStorage() >= size {

			lru.CreateNode(key, value)
			return true
		} else {
			for lru.RemainingStorage() < size {
				lru.DeleteNode(lru.back)
			}

			lru.CreateNode(key, value)
			return true
		}
	}
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.numBindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
