package cache

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	limit       int
	inUse       int
	numBindings int
	// key string, val: {starting point in array, # bytes}
	location map[string]*Node
	front    *Node
	back     *Node
	hits     int
	misses   int
	// current int
}

type Node struct {
	next     *Node
	previous *Node
	key      string
	value    []byte
	size     int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := new(FIFO)
	fifo.front = new(Node)
	fifo.back = fifo.front
	fifo.limit = limit
	fifo.location = make(map[string]*Node)
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
func (fifo *FIFO) PlaceNodeFront(current *Node) bool {
	if current == fifo.front {
		// Do nothing
	} else if current == fifo.back {
		// current=back; current -> front; current = front; front -> next
		fifo.back = current.previous
		fifo.back.next = nil

		current.next = fifo.front
		current.previous = nil
		fifo.front = current
	} else {
		// current=middle; current -> front; current = front; front -> next
		current.next.previous = current.previous
		current.previous.next = current.next
		current.next = fifo.front
		current.previous = nil
		fifo.front = current
	}
	return true
}
func (fifo *FIFO) DeleteNode(current *Node) bool {

	if fifo.numBindings == 1 {
	} else if current == fifo.back {
		// current=back; current -> front; current = front; front -> next
		fifo.back = current.previous
	} else if current == fifo.front {
		fifo.front = current.next
	} else {
		// current=middle; current -> front; current = front; front -> next
		current.next.previous = current.previous
		current.previous.next = current.next
	}

	fifo.inUse += -len(current.key) - len(fifo.location[current.key].value)
	delete(fifo.location, current.key)
	current = nil
	fifo.numBindings--

	return true
}
func (fifo *FIFO) CreateNode(key string, value []byte) bool {
	size := len(key) + len(value)

	if fifo.numBindings == 0 {
		fifo.front.key = key
		fifo.front.value = value
		fifo.location[key] = fifo.front

	} else {
		new_node := new(Node)
		new_node.next = fifo.front
		new_node.previous = nil
		new_node.value = value
		new_node.key = key
		new_node.size = size
		fifo.front.previous = new_node
		fifo.front = new_node

		fifo.location[key] = new_node
	}

	fifo.inUse += size
	fifo.numBindings++
	return true

}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	val, ok := fifo.location[key]
	if ok {
		// fifo.PlaceNodeFront(val)
		fifo.hits++
		return val.value, ok

	} else {
		fifo.misses++
		return nil, false
	}
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	if val, ok := fifo.location[key]; ok {
		fifo.DeleteNode(val)
		fifo.hits++
		return val.value, ok
	} else {
		fifo.misses++
		return nil, false

	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {

	size := len(key) + len(value)

	// if binding exists update old value, if size permits?????
	//  NOT SURE
	if val, ok := fifo.location[key]; ok {
		if val.size >= size {
			fifo.DeleteNode(val)
			fifo.CreateNode(key, value)
			return true
		}
		return false

		// else if a existing binding is not found then need to create from scratch
	} else {
		if size > fifo.limit {
			return false
		}

		if fifo.RemainingStorage() >= size {
			fifo.CreateNode(key, value)
			return true
		} else {
			for fifo.RemainingStorage() < size {
				fifo.DeleteNode(fifo.back)
			}
			fifo.CreateNode(key, value)
			return true
		}
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
