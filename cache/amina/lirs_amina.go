package cache

//package hashtable

import (
	"github.com/gammazero/deque"
)

const HIRS = 1
const HIRS_NR = 0
const LIRS_P = -1

type Element struct {
	// key
	key string
	// page -- if HIRS_NR, then == NIL
	page []byte
	// resident HIRS, non resident HIRS, or LIRS
	status int
}

type Location struct {
	key    string
	inS    bool
	indexS int
	inQ    bool
	indexQ bool
}

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LIRS struct {
	// max # of pages
	capacity int
	// capacity of LIR stack
	LIRcap int
	// capacity of HIR stack
	HIRcap int
	// # of pages in use
	inUse int
	//  LIRS stack
	S deque.Deque
	// resident HIR page stacks
	Q deque.Deque
	// tracks hits and misses
	stats *Stats
	// find location of element in queue S/Q
	location map[string]Location
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLirs(capacity int) *LIRS {
	lirs := new(LIRS)
	lirs.capacity = capacity
	lirs.inUse = 0
	//lirs.Q = deque.Deque{}
	//lirs.S = deque.Deque{}
	lirs.Q = *deque.New()
	lirs.S = *deque.New()
	lirs.location = make(map[string]Location)

	// init stats
	stats := &Stats{}
	stats.Hits = 0
	stats.Misses = 0
	lirs.stats = stats

	return lirs
}

// MaxStorage returns the maximum number of pages this LIRS can store
func (lirs *LIRS) MaxStorage() int {
	return lirs.capacity
}

// RemainingStorage returns the number of unused pages available in this LIRS
func (lirs *LIRS) RemainingStorage() int {
	return lirs.capacity - lirs.inUse
}

func PruneStack(S Deque.deque) (S Deque.deque) {

}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lirs *LIRS) Get(key string) (value []byte, ok bool) {
	Q := lirs.Q
	S := lirs.S

	foundKey := false
	for i := 0; i < S.Len(); i++ {
		elem, ok := S.At(i).(Element)
		if ok {
			if key == elem.key {
				foundKey = true

				// Upon accessing an LIR block X: This access is
				// guaranteed to be a hit in the cache. We move it to the
				// top of stack S. If the LIR block is originally located in
				// the bottom of the stack, we conduct a stack pruning.
				if elem.status == LIRS_P {
					lirs.stats.Hits += 1
					_ = S.Remove(i)
					S.PushFront(elem)
					lirs.S = PruneStack(S)
					return elem.page, true

					// always a miss. do not do anything.
				} else if elem.status == HIRS_NR {
					lirs.stats.Misses += 1
					return nil, false

					// Upon accessing an LIR block X: This access is
					// guaranteed to be a hit in the cache. We move it to the
					// top of stack S. If the LIR block is originally located in
					// the bottom of the stack, we conduct a stack pruning
				} else if elem.status == HIRS {
					lirs.stats.Hits += 1
					_ = S.Remove(i)

				}
			}

		} else {
			// error
			return nil, false
		}
	}
	if !foundKey {
		for i := 0; i < Q.Len(); i++ {
			elem, ok := S.At(i).(Element)
			if ok {
				if key == elem.key {
					foundKey = true
					if elem.status == HIRS {

					} else if elem.status == HIRS_NR {

						// elem.status = LIRS
					} else {

					}
				}
			} else {
				// error
				return nil, false

			}

		}
	}

	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lirs *LIRS) Remove(key string) (value []byte, ok bool) {
	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lirs *LIRS) Set(key string, value []byte) bool {
	Q := lirs.Q
	S := lirs.S

	foundKey := false
	for i := 0; i < S.Len(); i++ {
		elem, ok := S.At(i).(Element)
		if ok {
			if key == elem.key {
				foundKey = true
				_ = S.Remove(i)
				if elem.status == HIRS {

				} else if elem.status == HIRS_NR {

					// elem.status = LIRS
				} else {

				}
			}

		} else {
			// error
		}
	}
	if !foundKey {
		for i := 0; i < Q.Len(); i++ {
			elem, ok := S.At(i).(Element)
			if ok {
				if key == elem.key {
					foundKey = true
					if elem.status == HIRS {

					} else if elem.status == HIRS_NR {

						// elem.status = LIRS
					} else {

					}
				}
			} else {
				// error

			}

		}
	}
	if !foundKey {
		if lirs.capacity == lirs.inUse {

		} else {

			lirs.inUse += 1
		}
	}
	return false
}

// Len returns the number of bindings in the LIRS.
func (lirs *LIRS) Len() int {
	return lirs.inUse
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lirs *LIRS) Stats() *Stats {
	return &Stats{}
}
