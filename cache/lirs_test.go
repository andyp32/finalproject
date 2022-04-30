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

	//  for i := 0; i < 4; i++ {
	// 	 key := fmt.Sprintf("key%d", i)
	// 	 val := []byte(key)
	// 	 ok := fifo.Set(key, val)
	// 	 if !ok {
	// 		 t.Errorf("Failed to add binding with key: %s", key)
	// 		 t.FailNow()
	// 	 }
	// 	 // fmt.Println("Round", i)

	// 	 res, _ := fifo.Get(key)
	// 	 if !bytesEqual(res, val) {
	// 		 t.Errorf("Wrong value %s for binding with key: %s", res, key)
	// 		 t.FailNow()
	// 	 }
	//  }
}
