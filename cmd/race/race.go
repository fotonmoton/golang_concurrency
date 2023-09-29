package main

import (
	"fmt"
	"sync"
)

/*
1.  3---- 1---- 2---- 4----
2.
3.
4.

*/

// go run -race race.go
func main() {

	counter := 0

	const num = 1000
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			temp := counter
			// needed only for small numbers of gorountines
			// runtime.Gosched()
			temp++
			counter = temp
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("count:", counter)
}
