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
	"fmt"
	"testing"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

/******************************************************************************/
// test prune stacks
/******************************************************************************/

// should not remove anything of the end of stack S
func TestNoTailTrim(t *testing.T) {
	capacity := 10
	lirs := NewLIRS(capacity)
	key := "A"
	val := []byte(key)
	lirs.Set(key, val)

	key = "B"
	lirs.Set(key, val)

	key = "C"
	lirs.Set(key, val)
	lirs.GraphStacks()
	S := lirs.PruneStack(lirs.S, lirs.Q)
	lirs.S = S
	lirs.GraphStacks()

	// key = "D"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "E"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "F"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "G"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "H"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "I"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

	// key = "J"
	// lirs.Set(key, val)
	// lirs.GraphStacks()

}

// should remove X elements off the end of stack S
func TestTailTrim(t *testing.T) {
}

/******************************************************************************/
// test find
/******************************************************************************/

// find 0th element
func TestFindFirst(t *testing.T) {
}

// find last element
func TestFindLast(t *testing.T) {
}

// find element in middle of stack
func TestFindCenter(t *testing.T) {
}

// search for element not in stack
func TestFindNotIn(t *testing.T) {
}

/******************************************************************************/
// test insertion order into empty cache
/******************************************************************************/

func TestInsertEmpty3(t *testing.T) {
	capacity := 3
	lirs := NewLIRS(capacity)
	//checkCapacity(t, lirs, capacity)
	key := "E"
	val := []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test1")

	key = "A"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test2")

	key = "D"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test3")

	key = "B"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test4")

	lirs.Get("E")
	fmt.Println("test5")
	lirs.Remove("E")
	fmt.Println("test5")

	lirs.GraphStacks()
}

func TInsertEmpty5(t *testing.T) {
	capacity := 5
	lirs := NewLIRS(capacity)
	//checkCapacity(t, lirs, capacity)
	key := "E"
	val := []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test1")

	key = "A"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test2")

	key = "D"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test3")

	key = "B"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test4")

	lirs.Get("E")
	fmt.Println("test5")
	lirs.Remove("E")
	fmt.Println("test6")

	lirs.GraphStacks()
}

func TInsertEmpty10(t *testing.T) {
	capacity := 10
	lirs := NewLIRS(capacity)
	//checkCapacity(t, lirs, capacity)
	key := "E"
	val := []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test1")

	key = "A"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test2")

	key = "D"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test3")

	key = "B"
	val = []byte(key)
	lirs.Set(key, val)
	lirs.GraphStacks()
	fmt.Println("test4")

	lirs.Get("E")
	fmt.Println("test5")

	lirs.GraphStacks()
}
