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
	"testing"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

func TestLIRS(t *testing.T) {
	capacity := 64
	lirs := NewLIRS(capacity)
	checkCapacity(t, lirs, capacity)
	key := "E"
	val := []byte(key)
	lirs.Set(key, val, HIRS)

	key = "A"
	val = []byte(key)
	lirs.Set(key, val, LIRS_P)

	key = "D"
	val = []byte(key)
	lirs.Set(key, val, HIRS)
	key = "B"
	val = []byte(key)
	lirs.Set(key, val, LIRS_P)
	lirs.GraphStacks()

	lirs.Get("E")

	lirs.GraphStacks()
}
