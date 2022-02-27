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
	// capacity := 64
	// lru := NewLru(capacity)
	// checkCapacity(t, lru, capacity)
	// for i := 0; i < 4; i++ {

	// 	key := fmt.Sprintf("key%d", i)
	// 	val := []byte(key)
	// 	ok := lru.Set(key, val)

	// 	if !ok {
	// 		t.Errorf("Failed to add binding with key: %s", key)
	// 		t.FailNow()
	// 	}

	// 	res, _ := lru.Get(key)
	// 	if !bytesEqual(res, val) {
	// 		t.Errorf("Wrong value %s for binding with key: %s", res, key)
	// 		t.FailNow()
	// 	}

	// }
	// testkey := "key0"

	// res1, _ := lru.Remove(testkey)
	// val := []byte(testkey)

	// if !bytesEqual(res1, val) {
	// 	t.Errorf("Wrong value %s for removal of key: %s", res1, testkey)
	// 	t.FailNow()
	// }
	// res1, _ := lru.(testkey)
	// val := []byte(testkey)

	// if !bytesEqual(res1, val) {
	// 	t.Errorf("Wrong value %s for removal of key: %s", res1, testkey)
	// 	t.FailNow()
	// }
}
