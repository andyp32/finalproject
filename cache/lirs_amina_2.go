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
	key string
	inS bool
	//indexS int
	inQ bool
	//indexQ bool
	// resident HIRS, non resident HIRS, or LIRS
	status int
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

func (lirs *LIRS) PruneStack(S deque.Deque, Q deque.Deque) (S_ deque.Deque) {
	for true {
		elem, ok := S.Back().(Element)
		if ok {
			if elem.status == HIRS || elem.status == HIRS_NR {
				_ = S.PopBack()
				i, _ := Find(elem.key, Q)
				if i == -1 {
					// remove from location map
					delete(lirs.location, elem.key)
				}
			} else {
				return S
			}

		} else {
			return S
		}

	}
	return S
}

// if not in queue, return i = -1
func Find(key string, S deque.Deque) (i int, elem Element) {
	for i := 0; i < S.Len(); i++ {
		elem, ok := S.At(i).(Element)
		if ok {
			if elem.key == key {
				return i, elem
			}
		}
	}

	throwaway := new(Element)
	return -1, *throwaway

}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lirs *LIRS) Get(key string) (value []byte, ok bool) {
	Q := lirs.Q
	S := lirs.S

	location, ok := lirs.location[key]
	if ok {
		// Upon accessing an HIR resident block X: This
		// is a hit in the cache. We move it to the top of stack
		// S. There are two cases for block X: (1) If X is in the
		// stack S, we change its status to LIR. This block is also
		// removed from list Q. The LIR block in the bottom of
		// S is moved to the end of list Q with its status changed
		// to HIR. A stack pruning is then conducted. (2) If X
		// is not in stack S, we leave its status in HIR and move
		// it to the end of list Q.
		if location.status == HIRS {
			i, elem := Find(key, S)
			// not in stack S
			if i == -1 {
				// we leave its status in HIR and move it to the end of list Q
				j, _ := Find(key, Q)
				elem, ok = Q.Remove(j).(Element)
				if ok {
					Q.PushBack(elem)
					lirs.stats.Hits += 1
					lirs.S = S
					lirs.Q = Q
					return elem.page, true

				} else {
					//error
					return nil, false
				}

				// in stack S
			} else {
				// change its status to LIR
				elem.status = LIRS_P
				location.status = LIRS_P

				// move to top of stack S
				_ = S.Remove(i)
				S.PushFront(elem)

				// remove from list Q
				j, _ := Find(key, Q)
				_ = Q.Remove(j)

				// The LIR block in the bottom of S is moved to the end of list Q with its status changed to HIR.
				elemBack, ok := S.Back().(Element)
				if ok {
					if elemBack.status == LIRS_P {
						elemBack.status = HIRS
						_ = S.PopBack()
						Q.PushBack(elemBack)
					}
				}

				// A stack pruning is then conducted
				S = lirs.PruneStack(S, Q)

				// update stats and return page
				lirs.location[key] = location
				lirs.S = S
				lirs.Q = Q
				lirs.stats.Hits += 1
				return elem.page, true
			}
			// Upon accessing an LIR block X: This access is
			// guaranteed to be a hit in the cache. We move it to the
			// top of stack S. If the LIR block is originally located in
			// the bottom of the stack, we conduct a stack pruning.
		} else if location.status == LIRS_P {
			i, elem := Find(key, S)
			// not in stack S
			if i == -1 {
				// error
				return nil, false
				// in stack S
			} else {
				// move to top of stack S
				_ = S.Remove(i)
				S.PushFront(elem)
				S = lirs.PruneStack(S, Q)

				// update stats and return page
				lirs.S = S
				lirs.Q = Q
				lirs.stats.Hits += 1
				return elem.page, true
			}

			// default miss
			// if location.status == HIRS_NR
		} else {
			lirs.stats.Misses += 1
			return nil, false
		}
		// not in cache
	} else {
		lirs.stats.Misses += 1
		return nil, false
	}
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lirs *LIRS) Remove(key string) (value []byte, ok bool) {
	Q := lirs.Q
	S := lirs.S

	ok = false
	i, elem1 := Find(key, S)

	if i != -1 {
		S.Remove(i)
		S = lirs.PruneStack(S, Q)
		delete(lirs.location, key)
		ok = true

	}

	j, elem2 := Find(key, Q)
	if j != -1 {
		Q.Remove(i)
		ok = true

	}

	lirs.S = S
	lirs.Q = Q

	if i != -1 {
		return elem1.page, true
	} else if j != -1 {
		return elem2.page, true
	} else {
		return nil, false
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lirs *LIRS) Set(key string, value []byte) bool {
	Q := lirs.Q
	S := lirs.S

	location, ok := lirs.location[key]
	if ok {
		// Upon accessing an HIR resident block X: This
		// is a hit in the cache. We move it to the top of stack
		// S. There are two cases for block X: (1) If X is in the
		// stack S, we change its status to LIR. This block is also
		// removed from list Q. The LIR block in the bottom of
		// S is moved to the end of list Q with its status changed
		// to HIR. A stack pruning is then conducted. (2) If X
		// is not in stack S, we leave its status in HIR and move
		// it to the end of list Q.
		if location.status == HIRS {
			i, elem := Find(key, S)
			// not in stack S
			if i == -1 {
				// we leave its status in HIR and move it to the end of list Q
				j, _ := Find(key, Q)
				elem, ok = Q.Remove(j).(Element)
				if ok {
					elem.page = value
					Q.PushBack(elem)
					lirs.stats.Hits += 1
					lirs.S = S
					lirs.Q = Q

				} else {
					//error
					return false
				}

				// in stack S
			} else {
				// change its status to LIR
				elem.status = LIRS_P
				location.status = LIRS_P

				// move to top of stack S
				_ = S.Remove(i)
				elem.page = value
				S.PushFront(elem)

				// remove from list Q
				j, _ := Find(key, Q)
				_ = Q.Remove(j)

				// The LIR block in the bottom of S is moved to the end of list Q with its status changed to HIR.
				elemBack, ok := S.Back().(Element)
				if ok {
					if elemBack.status == LIRS_P {
						elemBack.status = HIRS
						_ = S.PopBack()
						Q.PushBack(elemBack)
					}
				}

				// A stack pruning is then conducted
				S = lirs.PruneStack(S, Q)

				// update stats and return page
				lirs.location[key] = location
				lirs.S = S
				lirs.Q = Q
				lirs.stats.Hits += 1
			}
			// Upon accessing an LIR block X: This access is
			// guaranteed to be a hit in the cache. We move it to the
			// top of stack S. If the LIR block is originally located in
			// the bottom of the stack, we conduct a stack pruning.
		} else if location.status == LIRS_P {
			i, elem := Find(key, S)
			// not in stack S
			if i == -1 {
				// error
				// in stack S
			} else {
				// move to top of stack S
				_ = S.Remove(i)
				elem.page = value
				S.PushFront(elem)
				S = lirs.PruneStack(S, Q)

				// update stats and return page
				lirs.S = S
				lirs.Q = Q
				lirs.stats.Hits += 1
			}

			// Upon accessing an HIR non-resident block X:
			// This is a miss. We remove the HIR resident block
			// at the front of list Q (it then becomes a non-resident
			// block), and replace it out of the cache. Then we load
			// the requested block X into the freed buffer and place
			// it on the top of stack S. There are two cases for block
			// X: (1) If X is in stack S, we change its status to LIR
			// and move the LIR block in the bottom of stack S to
			// the end of list Q with its status changed to HIR. A
			// stack pruning is then conducted. (2) If X is not in
			// stack S, we leave its status in HIR and place it in the
			// end of list Q.

			// if location.status == HIRS_NR
		} else {

			// cache is full
			if Q.Len() > 0 && lirs.inUse == lirs.capacity {
				// remove the HIR resident block at the front of list Q
				elem, ok := Q.PopFront().(Element)
				if ok {
					// update elem to HIRS_NR in S
					i, _ := Find(elem.key, S)
					if i != -1 {
						elem.status = HIRS_NR
						S.Set(i, elem)
						// elem, ok := S.Remove(i).(Element)
						// if ok {
						// 	elem.status = HIRS_NR
						// 	S.PushFront(elem)
						// }

						// update location in map
						location = lirs.location[elem.key]
						location.status = HIRS_NR
						lirs.location[key] = location
					} else {
						// remove from location map
						delete(lirs.location, elem.key)
					}
					// cache is not full
				}
			} else if lirs.inUse < lirs.capacity {
				lirs.inUse += 1
				// cache is full but Q is empty -- error!
			} else {
				return false

			}

			// X is in stack S, we change its status to LIR
			i, _ := Find(key, S)
			if i != -1 {
				elem, ok := S.Remove(i).(Element)
				if ok {
					elem.status = LIRS_P
					S.PushFront(elem)

					// update location in map
					location = lirs.location[elem.key]
					location.status = HIRS_NR
					lirs.location[key] = location

					// and move the LIR block in the bottom of stack S to
					// the end of list Q with its status changed to HIR.
					elem, ok = S.PopBack().(Element)
					if ok {
						elem.status = HIRS
						Q.PushBack(elem)
						S = lirs.PruneStack(S, Q)

						// update location in map
						location = lirs.location[elem.key]
						location.status = HIRS
						lirs.location[key] = location
					}

				}

			} else {
				// error
				return false
			}

		}

		lirs.stats.Misses += 1
		// not in cache
	} else {
		// cache is full
		if Q.Len() > 0 && lirs.inUse == lirs.capacity {
			// remove the HIR resident block at the front of list Q
			elem, ok := Q.PopFront().(Element)
			if ok {
				// update elem to HIRS_NR in S
				i, _ := Find(elem.key, S)
				if i != -1 {
					elem.status = HIRS_NR
					S.Set(i, elem)
					// elem, ok := S.Remove(i).(Element)
					// if ok {
					// 	elem.status = HIRS_NR
					// 	S.PushFront(elem)
					// }

					// update location in map
					location = lirs.location[elem.key]
					location.status = HIRS_NR
					lirs.location[key] = location
				} else {
					// remove from location map
					delete(lirs.location, elem.key)
				}
				// cache is not full
			}
		} else if lirs.inUse < lirs.capacity {
			lirs.inUse += 1
			// cache is full but Q is empty -- error!
		} else {
			return false

		}

		// If X is not in
		// stack S, we leave its status in HIR and place it in the
		// end of list Q.
		elem := new(Element)
		elem.status = HIRS
		elem.key = key
		elem.page = value
		S.PushBack(elem)
		Q.PushBack(elem)

		// update location in map
		location = lirs.location[elem.key]
		location.status = HIRS
		lirs.location[key] = location

		lirs.stats.Misses += 1
	}

	lirs.Q = Q
	lirs.S = S
	return true
}

// // Set associates the given value with the given key, possibly evicting values
// // to make room. Returns true if the binding was added successfully, else false.
// func (lirs *LIRS) Set(key string, value []byte) bool {
// 	Q := lirs.Q
// 	S := lirs.S

// 	foundKey := false
// 	for i := 0; i < S.Len(); i++ {
// 		elem, ok := S.At(i).(Element)
// 		if ok {
// 			if key == elem.key {
// 				foundKey = true
// 				_ = S.Remove(i)
// 				if elem.status == HIRS {

// 				} else if elem.status == HIRS_NR {

// 					// elem.status = LIRS
// 				} else {

// 				}
// 			}

// 		} else {
// 			// error
// 		}
// 	}
// 	if !foundKey {
// 		for i := 0; i < Q.Len(); i++ {
// 			elem, ok := S.At(i).(Element)
// 			if ok {
// 				if key == elem.key {
// 					foundKey = true
// 					if elem.status == HIRS {

// 					} else if elem.status == HIRS_NR {

// 						// elem.status = LIRS
// 					} else {

// 					}
// 				}
// 			} else {
// 				// error

// 			}

// 		}
// 	}
// 	if !foundKey {
// 		if lirs.capacity == lirs.inUse {

// 		} else {

// 			lirs.inUse += 1
// 		}
// 	}
// 	return false
// }

// Len returns the number of bindings in the LIRS.
func (lirs *LIRS) Len() int {
	return lirs.inUse
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lirs *LIRS) Stats() *Stats {
	return &Stats{}
}
