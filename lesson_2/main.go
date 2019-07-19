package main

import (
	"fmt"
	"time"
)



func main() {
	intCh0 := make(chan int, 10)
	intCh1 := make(chan int, 10)

	go factorial(5, intCh0)

	//for	k := range intCh0 {
	//	fmt.Println(k)
	//}
	L:
	for {
		select {
		case b := <-intCh0:
			if b == 0 {
				break L
			}
			fmt.Printf("channel intCh0 : %v\n", b)
		case v := <-intCh1:
			fmt.Printf("channel intCh1: %v\n", v)
		case <-time.After(5 * time.Second):
			fmt.Printf("timeout")
			//default:
			//	fmt.Printf("channel zero: %v", 0)
		}
	}


}

func factorial(n int, ch chan int) {
	defer close(ch)
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
		// fmt.Printf("%v\n", result)
		ch <- result // посылаем по числу
	}
}
