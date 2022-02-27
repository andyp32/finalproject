/******************************************************************************
 * lru_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for lru.go. You are welcome to change
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

func TestLRU(t *testing.T) {
	capacity := 100
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)
	for i := 0; i < 5; i++ {
		for j := 0; j < 2; j++ {

			key := fmt.Sprintf("___%d%d", i, j)
			val := []byte(key)
			ok := lru.Set(key, val)

			if !ok {
				t.Errorf("Failed to add binding with key: %s", key)
				t.FailNow()
			}

			// res, _ := lru.Get(key)
			// if !bytesEqual(res, val) {
			// 	t.Errorf("Wrong value %s for binding with key: %s", res, key)
			// 	t.FailNow()
			// }
		}
	}
	for i := 3; i < 5; i++ {
		for j := 0; j < 2; j++ {

			key := fmt.Sprintf("___%d%d", i, j)
			val := []byte(key)
			// ok := lru.Set(key, val)

			// if !ok {
			// 	t.Errorf("Failed to add binding with key: %s", key)
			// 	t.FailNow()
			// }

			res, _ := lru.Get(key)
			if !bytesEqual(res, val) {
				t.Errorf("Wrong value %s for binding with key: %s", res, key)
				t.FailNow()
			}
		}
	}

}
