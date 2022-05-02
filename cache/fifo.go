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
	previous  *Node
	key       string
	lirs_type int
	next      *Node
	value     []byte
	size      int
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

func (fifo *FIFO) DeleteNode(current *Node) bool {

	if fifo.numBindings == 1 {
	} else if current == fifo.back {
		// current=back; current -> front; current = front; front -> next
		fifo.back = current.previous
		fifo.back.next = nil
		// fmt.Println("here")
	} else if current == fifo.front {
		fifo.front = current.next
		fifo.front.previous = nil
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
func (fifo *FIFO) CreateNode(key string, value []byte, lirs_type int) bool {
	size := len(key) + len(value)

	if fifo.numBindings == 0 {
		fifo.front.key = key
		fifo.front.value = value
		fifo.front.lirs_type = lirs_type
		fifo.location[key] = fifo.front
	} else {
		new_node := new(Node)
		new_node.previous = fifo.back
		new_node.next = nil
		new_node.value = value
		new_node.key = key
		new_node.size = size
		new_node.lirs_type = lirs_type
		fifo.back.next = new_node
		fifo.back = new_node

		fifo.location[key] = new_node
	}

	fifo.inUse += size
	fifo.numBindings++
	return true

}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value *Node, ok bool) {
	val, ok := fifo.location[key]
	if ok {
		fifo.hits++
		return val, ok
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
		return val.value, ok
	} else {
		return nil, false
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte, lirs_type int) bool {

	size := len(key) + len(value)
	if size > fifo.limit {
		return false
	}
	// if binding exists update old value, if size permits
	if val, ok := fifo.location[key]; ok {
		bytes_diff := len(value) - len(val.value) // check space needed
		if fifo.RemainingStorage() >= bytes_diff {
			val.value = value
			fifo.inUse += bytes_diff
			val.size = size
			val.lirs_type = lirs_type
			return true
		}
		return false

		// else if a existing binding is not found then need to create from scratch
	} else {
		if fifo.RemainingStorage() >= size {
			fifo.CreateNode(key, value, lirs_type)
			return true

		} else {
			for fifo.RemainingStorage() < size {
				fifo.DeleteNode(fifo.back)
			}
			fifo.CreateNode(key, value, lirs_type)
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
