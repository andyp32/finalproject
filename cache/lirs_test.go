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
func TestBasicGetSuccess(t *testing.T) {
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
func TestBasicGetFail(t *testing.T) {
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

// run a more complex series of queries, trace alg, and confirm GET works as expected
func TestComplexGetSuccess(t *testing.T) {
	lirs := MakeStack()

	// override "A"
	lirs.Set("F", []byte("F"))
	// lirs.GraphStacks()

	// should find "F"
	elem, _ := lirs.Get("F")
	if bytes.Compare(elem, []byte("F")) != 0 {
		t.Errorf("GET failed")
	}

	// remove "F"
	lirs.Remove("F")
	// lirs.GraphStacks()

	// insert "G"
	lirs.Set("G", []byte("G"))
	// lirs.GraphStacks()

	// should find "G"
	elem, _ = lirs.Get("G")
	if bytes.Compare(elem, []byte("G")) != 0 {
		t.Errorf("GET failed")
	}

}

// run a more complex series of queries, trace alg, and confirm GET works as expected
func TestComplexGetFail(t *testing.T) {
	lirs := MakeStack()

	// override "A"
	lirs.Set("F", []byte("F"))
	// lirs.GraphStacks()

	// should not find "A"
	_, ok := lirs.Get("A")
	if ok == true {
		t.Errorf("GET failed")
	}

	// remove "F"
	lirs.Remove("F")
	// lirs.GraphStacks()

	// // should not find "F"
	// _, ok = lirs.Get("F")
	// if ok == true {
	// 	t.Errorf("GET failed")
	// }
}

// /******************************************************************************/
// // test REMOVE
// /******************************************************************************/

// remove elements in the cache
func TestRemoveBasicSuccess(t *testing.T) {

}

// attempt to remove elements not in the cache
func TestRemoveBasicFailure(t *testing.T) {

}

// run a more complex series of queries, trace alg, and confirm REMOVE works as expected
func TestRemoveComplexSuccess(t *testing.T) {

}

// run a more complex series of queries, trace alg, and confirm REMOVE works as expected
func TestRemoveComplexFailure(t *testing.T) {

}

// /******************************************************************************/
// // test SET
// /******************************************************************************/

// remove elements in the cache
func TestSetBasicSuccess(t *testing.T) {

}

// attempt to remove elements not in the cache
func TestSetBasicFailure(t *testing.T) {

}

// run a more complex series of queries, trace alg, and confirm SET works as expected
func TestSetComplexSuccess(t *testing.T) {

}

// run a more complex series of queries, trace alg, and confirm SET works as expected
func TestSetComplexFailure(t *testing.T) {

}

// /******************************************************************************/
// test COMPLEX COMBINATIONS
// /******************************************************************************/

// test sequences of REMOVE, GET, SET
// check cache hits/misses as evidence of correctness
func TestAlg1(t *testing.T) {

}

func TestAlg2(t *testing.T) {

}

func TestAlg3(t *testing.T) {

}

func TestAlg4(t *testing.T) {

}

func TestAlg5(t *testing.T) {

}

func TestAlg6(t *testing.T) {

}

func TestAlg7(t *testing.T) {

}

func TestAlg8(t *testing.T) {

}

func TestAlg9(t *testing.T) {

}

func TestAlg10(t *testing.T) {

}
