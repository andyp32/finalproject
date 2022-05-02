package cache

import "fmt"

const HIRS = 1
const HIRS_NR = 0
const LIRS_P = -1

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

func (lirs *LIRS) PlaceNodeFront(stack *FIFO, current *Node) bool {
	if current == stack.front {
		// Do nothing
	} else if stack.numBindings == 2 {
		// current=back; current -> front; current = front; front -> next
		stack.back = current.previous

		current.next = stack.front
		stack.front = current
		stack.front.next = stack.back
		stack.back.previous = stack.front
		stack.front.previous = nil
		stack.back.next = nil

	} else if current == stack.back {
		// current=back; current -> front; current = front; front -> next
		stack.back = current.previous
		current.next = stack.front
		stack.front.previous = current
		stack.front = current

		stack.front.previous = nil
		stack.back.next = nil

	} else {
		// current=middle; current -> front; current = front; front -> next
		current.next.previous = current.previous
		current.previous.next = current.next

		current.next = stack.front
		stack.front.previous = current
		stack.front = current
		stack.front.previous = nil
	}
	return true
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (lirs *LIRS) Get(key string) (value *Node, ok bool) {

	node, ok := lirs.S.Get(key)
	if ok {
		lirs.hits++
	} else {
		lirs.misses++
		// miss and a resident page has to be replaced, the resident HIR page at the
		// bottom of Stack Q is selected as the victim for replacement
		// Q_back := lirs.Q.back
		// lirs.Q.Remove(Q_back.key)
		// lirs.Q.Set()
		return nil, false
	}
	if node.lirs_type == HIRS {
		node.lirs_type = LIRS_P
	}

	// lirs.PlaceNodeFront(lirs.S, node)
	current := lirs.S.back

	// LIR page currently at Stack S’s bottom turns into a HIR page and moves to the top of Stack Q
	if current.lirs_type == LIRS_P {
		lirs.Q.Set(current.key, current.value, HIRS)
		lirs.S.Remove(current.key)
	}

	// Remove any HIRS at the bottom of stack
	current = lirs.S.back
	for current != nil {
		if current.lirs_type != HIRS {
			break
		}
		key := current.key
		current = current.previous
		lirs.S.Remove(key)
	}

	return
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lirs *LIRS) Set(key string, value []byte, lirs_type int) bool {

	lirs.S.Set(key, value, lirs_type)
	if lirs_type == HIRS {
		lirs.Q.Set(key, value, lirs_type)
	}
	return true
}
func (lirs *LIRS) GraphStacks() {
	fmt.Println("========================================")

	current := lirs.S.front
	for current != nil {
		fmt.Print(" -> ")
		fmt.Print(current.key, current.lirs_type)
		current = current.next
	}
	fmt.Println()
	current = lirs.Q.front
	for current != nil {
		fmt.Print(" -> ")
		fmt.Print(current.key, current.lirs_type)
		current = current.next
	}
	fmt.Println()

	fmt.Println("========================================")

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
