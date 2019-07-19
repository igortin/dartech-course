package main

import "fmt"

func main() {

	intCh := make(chan int,5)
	intCh <- 9
	intCh <- 9; intCh <- 999
	close(intCh) // канал закрыт
	for {
		if val, opened := <-intCh; opened {
			fmt.Println(val)
		} else {
			break
			}
	}}