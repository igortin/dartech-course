package main

import (
"fmt"
"os"
"strconv"
"sync"
)

var (
	path  = "./output.txt"
	count = 12
)

func main() {
	var wg sync.WaitGroup
	var mnOnce sync.Once
	wg.Add(5)
	ch := make(chan []int)
	ch1 := make(chan []int)
	ch2 := make(chan []int)

	defer close(ch)
	defer close(ch1)
	defer close(ch2)

	go func(n int, ch chan<- []int) {
		var t []int
		for k := 0; k < count; k++ {
			t = append(t, fibonacci(k))
		}
		ch <- t
		wg.Done()
	}(count, ch)

	num := <-ch

	go func(num []int, ch1 chan []int) {
		ch1 <- num
		wg.Done()
	}(num, ch1)

	go func(num []int, ch2 chan []int) {
		ch2 <- num
		wg.Done()
	}(num, ch2)

	go func(ch1 chan []int) {
		fmt.Printf("%v", <-ch1)
		wg.Done()
	}(ch1)

	mnOnce.Do(func() {
		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			_, err := os.Create(path)

			if err != nil {
				os.Exit(1)
			}
		}
	})

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	go func(ch2 chan []int) {
		t := <-ch2
		for i := 0; i < len(t); i++ {
			b := []byte(strconv.Itoa(t[i])) // это корректно ?????
			_, _ = file.Write(b)
		}
		wg.Done()
	}(ch2)
	wg.Wait()
}

func fibonacci(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	if n == 0 {
		return 0
	}
	return fibonacci(n-1) + fibonacci(n-2)
}




