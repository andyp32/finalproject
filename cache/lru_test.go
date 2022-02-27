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
	capacity := 20
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)
	// for i := 0; i < 4; i++ {

	// 	key := fmt.Sprintf("____%d", i)
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
	key := fmt.Sprintf("1111")
	val := []byte("aaaa")
	lru.Set(key, val)
	key1 := fmt.Sprintf("2222")
	val1 := []byte("bbbb")
	lru.Set(key1, val1)
	// fmt.Println("BACK0000>>>>>>>>>", lru.back)

	// fmt.Println(lru.location[key1])
	a, _ := lru.Get(key)
	b, _ := lru.Get(key1)
	fmt.Println(a, b)

	key2 := fmt.Sprintf("3333")
	val2 := []byte("cccc")
	// fmt.Println("setting")
	// fmt.Println("BACK111111>>>>>>>>>", lru.back)

	lru.Set(key2, val2)
	// fmt.Println("done")
	// fmt.Println("BACK22222>>>>>>>>>", lru.back)

	c, _ := lru.Get(key)
	fmt.Println(c)
	// fmt.Println("BACK33333>>>>>>>>>", lru.back)

	lru.Remove("2222")
	// fmt.Println("BACK>>>>>>>>>", lru.back)

	lru.Get("2222")
	lru.Get("3333")

	fmt.Println(lru.Stats())
	// testkey := "____0"

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
