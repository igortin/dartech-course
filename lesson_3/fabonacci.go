package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	path  = "./output.txt"
	count = 29
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	ch := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func(n int, ch chan<- int) {
		defer close(ch)
		ch <- fibonacci(n)
		wg.Done()
	}(count, ch)

	num := <-ch

	go func(num int, ch1 chan int) {
		defer close(ch1)
		ch1 <- num
		wg.Done()
	}(num, ch1)

	go func(num int, ch2 chan int) {
		defer close(ch2)
		ch2 <- num
		wg.Done()
	}(num, ch2)

	go func(ch1 chan int) {
		fmt.Printf("%v", <-ch1)
		wg.Done()
	}(ch1)

	go func(ch2 chan int) {
		err := CheckFile(path)
		if err != nil {
			os.Exit(1)
		}
		file, err := os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		_, _ = file.WriteString(strconv.Itoa(<-ch2))
		wg.Done()
	}(ch2)

	wg.Wait()
}

func fibonacci(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	if n == 0 { return 0}
	return fibonacci(n-1) + fibonacci(n-2)
}

func CheckFile(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}
	return nil
}
