/******************************************************************************
 * fifo_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for fifo.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
 ******************************************************************************/

package cache

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gammazero/deque"
)

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

/******************************************************************************/
// HELPER
/******************************************************************************/

func PrintStats(stats *Stats) {
	fmt.Println("===========================")
	fmt.Println("Hits: ", stats.Hits)
	fmt.Println("Misses: ", stats.Misses)
}
func PrintStatus(i int) {
	if i == LIRS_P {
		fmt.Print("LIR")
	} else if i == HIRS {
		fmt.Print("HIR")

	} else {
		fmt.Print("HIR-NR")
	}
}

func (lirs *LIRS) GraphStacks() (S_ *deque.Deque, Q_ *deque.Deque) {
	fmt.Println("========================================")

	S := lirs.S
	Q := lirs.Q
	for i := 0; i < S.Len(); i++ {
		current, _ := S.At(i).(*Element)
		fmt.Print(" -> ")
		fmt.Print("(key: ", current.key, ", status: ")
		PrintStatus(current.status)
		fmt.Print(")")
	}
	fmt.Println()
	for i := 0; i < Q.Len(); i++ {
		current, _ := Q.At(i).(*Element)
		fmt.Print(" -> ")
		fmt.Print("(key: ", current.key, ", status: ")
		PrintStatus(current.status)
		fmt.Print(")")
	}
	fmt.Println()

	fmt.Println("========================================")

	return S, Q

}

func MakeStack() (lirs_ *LIRS) {
	capacity := 5
	lirs := NewLIRS(capacity)

	lirs.Set("A", []byte("A"))
	lirs.Set("B", []byte("B"))
	lirs.Set("C", []byte("C"))
	lirs.Set("D", []byte("D"))
	lirs.Set("E", []byte("E"))

	return lirs

}

/******************************************************************************/
// test PRUNESTACKS
/******************************************************************************/

// should not remove anything of the end of stack S
func TestNoTailNo(t *testing.T) {
	capacity := 10
	lirs := NewLIRS(capacity)
	val := []byte("X")
	lirs.Set("A", val)
	lirs.Set("B", val)
	lirs.Set("C", val)
	lirs.Set("D", val)
	lirs.Set("E", val)
	lirs.Set("F", val)
	lirs.Set("G", val)
	lirs.Set("H", val)
	lirs.Set("I", val)
	lirs.Set("J", val)
	lirs.Set("K", val)

	lirs.Get("A")
	lirs.Get("B")
	lirs.Get("D")
	lirs.Get("A")
	lirs.Get("F")
	lirs.Get("G")

	lirs.GraphStacks()
	if lirs.S.Len() != 1 {
		t.Errorf("Building cache failed")
	}

	lirs.PruneStack(lirs.S, lirs.Q)
	lirs.GraphStacks()

	if lirs.S.Len() != 1 {
		t.Errorf("Pruning failed")
	}

}

// should not 2 blocks at the end of stack S
func TestNoTailSome(t *testing.T) {
	capacity := 10
	lirs := NewLIRS(capacity)
	val := []byte("X")
	lirs.Set("A", val)
	lirs.Set("B", val)
	lirs.Set("C", val)
	lirs.Set("D", val)
	lirs.Set("E", val)
	lirs.Set("F", val)
	lirs.Set("G", val)
	lirs.Set("H", val)
	lirs.Set("I", val)
	lirs.Set("J", val)
	lirs.Set("K", val)

	lirs.Get("A")
	lirs.Get("B")
	lirs.Get("D")
	lirs.Get("A")
	lirs.Get("F")
	lirs.Get("G")
	lirs.Set("L", val)
	lirs.Set("M", val)
	lirs.GraphStacks()

	if lirs.S.Len() != 3 {
		t.Errorf("Building cache failed")
	}

	lirs.PruneStack(lirs.S, lirs.Q)
	lirs.GraphStacks()

	if lirs.S.Len() != 1 {
		t.Errorf("Pruning failed")
	}

}

// should remove all elements off the end of stack S
func TestTailTrim(t *testing.T) {
	lirs := MakeStack()

	lirs.PruneStack(lirs.S, lirs.Q)
	if lirs.S.Len() != 0 {
		t.Errorf("Pruning failed")
	}

}

// /******************************************************************************/
// // test FIND
// /******************************************************************************/

// find 0th element
func TestFindFirst(t *testing.T) {
	lirs := MakeStack()
	i, elem := Find("A", lirs.S)

	if i != 0 || elem.key != "A" {
		t.Errorf("Found wrong element")
	}

}

// find last element
func TestFindLast(t *testing.T) {
	lirs := MakeStack()
	i, elem := Find("E", lirs.S)

	if i != 4 || elem.key != "E" {
		t.Errorf("Found wrong element")
	}
}

// find element in middle of stack
func TestFindCenter(t *testing.T) {
	lirs := MakeStack()
	i, elem := Find("C", lirs.S)

	if i != 2 || elem.key != "C" {
		t.Errorf("Found wrong element")
	}
}

// search for element not in stack
func TestFindNotIn(t *testing.T) {

	lirs := MakeStack()
	i, _ := Find("F", lirs.S)

	if i != -1 {
		t.Errorf("Should not find any elements")
	}
}

// test find on the empty stack
func TestFindEmpty(t *testing.T) {
	capacity := 5
	lirs := NewLIRS(capacity)

	i, _ := Find("A", lirs.S)

	if i != -1 {
		t.Errorf("Should not find any elements")
	}
}

// /******************************************************************************/
// // test MAXSTORAGE and LEN functions
// /******************************************************************************/

func TestMaxStorage(t *testing.T) {
	capacity := 5
	lirs := NewLIRS(capacity)
	if capacity != lirs.capacity {
		t.Errorf("Capacity incorrect")
	}

	capacity = 10
	lirs = NewLIRS(capacity)
	if capacity != lirs.capacity {
		t.Errorf("Capacity incorrect")
	}
}

func TestLen(t *testing.T) {
	capacity := 5
	lirs := NewLIRS(capacity)
	val := []byte("X")

	key := "A"
	lirs.Set(key, val)

	if 1 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	key = "B"
	lirs.Set(key, val)

	if 2 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	key = "C"
	lirs.Set(key, val)

	if 3 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	key = "D"
	lirs.Set(key, val)

	if 4 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	key = "E"
	lirs.Set(key, val)

	if 5 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	// remove
	lirs.Remove("E")

	if 4 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	key = "E"
	lirs.Set(key, val)

	if 5 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	// replace existing key-value
	key = "E"
	lirs.Set(key, val)

	if 5 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

	// evict
	key = "F"
	lirs.Set(key, val)

	if 5 != lirs.Len() {
		t.Errorf("inUse incorrect")
	}

}

// /******************************************************************************/
// // test STATS
// /******************************************************************************/

// /******************************************************************************/
// // test GET
// /******************************************************************************/

// fill cache and query every element in cache
func TestBasicGet(t *testing.T) {
	lirs := MakeStack()
	value, _ := lirs.Get("A")
	if bytes.Compare(value, []byte("A")) != 0 {
		t.Errorf("GET failed")
	}

	value, _ = lirs.Get("B")
	if bytes.Compare(value, []byte("B")) != 0 {
		t.Errorf("GET failed")
	}

	value, _ = lirs.Get("C")
	if bytes.Compare(value, []byte("C")) != 0 {
		t.Errorf("GET failed")
	}

	value, _ = lirs.Get("D")
	if bytes.Compare(value, []byte("D")) != 0 {
		t.Errorf("GET failed")
	}

	value, _ = lirs.Get("E")
	if bytes.Compare(value, []byte("E")) != 0 {
		t.Errorf("GET failed")
	}

}

// query elements not in the cache
func TestBasicNot(t *testing.T) {
	lirs := MakeStack()
	_, ok := lirs.Get("F")
	if ok == true {
		t.Errorf("GET failed")
	}

	lirs.Remove("E")
	_, ok = lirs.Get("E")
	if ok == true {
		t.Errorf("GET failed")
	}

}

// /******************************************************************************/
// // test REMOVE
// /******************************************************************************/

// /******************************************************************************/
// // test SET
// /******************************************************************************/

// /******************************************************************************/
// // test insertion order into empty cache
// /******************************************************************************/

func TestAlg(t *testing.T) {
	capacity := 5
	lirs := NewLIRS(capacity)
	val := []byte("X")

	key := "A"
	lirs.Set(key, val)
	lirs.GraphStacks()

	key = "B"
	lirs.Set(key, val)
	lirs.GraphStacks()

	key = "C"
	lirs.Set(key, val)
	lirs.GraphStacks()

	key = "D"
	lirs.Set(key, val)
	lirs.GraphStacks()

	key = "E"
	lirs.Set(key, val)
	lirs.GraphStacks()

	lirs.Get("C")
	lirs.GraphStacks()

	lirs.Get("A")
	lirs.GraphStacks()

	lirs.Set("F", val)
	lirs.GraphStacks()

	lirs.Set("G", val)
	lirs.GraphStacks()

	fmt.Println(lirs.inUse)
	fmt.Println(lirs.capacity)

}

// func TestInsertEmpty3(t *testing.T) {
// 	// capacity := 3
// 	// lirs := NewLIRS(capacity)
// 	// //checkCapacity(t, lirs, capacity)
// 	// key := "E"
// 	// val := []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test1")

// 	// key = "A"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test2")

// 	// key = "D"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test3")

// 	// key = "B"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test4")

// 	// lirs.Get("E")
// 	// fmt.Println("test5")
// 	// lirs.Remove("E")
// 	// fmt.Println("test5")

// 	// lirs.GraphStacks()
// }

// func TInsertEmpty5(t *testing.T) {
// 	// capacity := 5
// 	// lirs := NewLIRS(capacity)
// 	// //checkCapacity(t, lirs, capacity)
// 	// key := "E"
// 	// val := []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test1")

// 	// key = "A"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test2")

// 	// key = "D"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test3")

// 	// key = "B"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test4")

// 	// lirs.Get("E")
// 	// fmt.Println("test5")
// 	// lirs.Remove("E")
// 	// fmt.Println("test6")

// 	// lirs.GraphStacks()
// }

// func TInsertEmpty10(t *testing.T) {
// 	// capacity := 10
// 	// lirs := NewLIRS(capacity)
// 	// //checkCapacity(t, lirs, capacity)
// 	// key := "E"
// 	// val := []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test1")

// 	// key = "A"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test2")

// 	// key = "D"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test3")

// 	// key = "B"
// 	// val = []byte(key)
// 	// lirs.Set(key, val)
// 	// lirs.GraphStacks()
// 	// fmt.Println("test4")

// 	// lirs.Get("E")
// 	// fmt.Println("test5")

// 	// lirs.GraphStacks()
// }
