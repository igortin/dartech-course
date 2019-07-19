package main

import (
	"fmt"
	"sync"
)

var a []int

func main() {
	var mu sync.Mutex
	var muOnce sync.Once
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {

		wg.Add(1)

		go func(k int) {
			defer wg.Done()
			mu.Lock()
			a = append(a, k)
			mu.Unlock()
		}(i)

		muOnce.Do(
			func() {
				fmt.Printf("only once %v\n", 555)
			})
		}

	muOnce.Do(
		func() {
			fmt.Printf("already once %v\n", 555)
		})


	wg.Wait()
	fmt.Println(a)
}