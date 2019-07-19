package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//var elements map[byte]int

	ch := make(chan map[byte]int)
	defer close(ch)
	res := make(map[string]int)

	slice := GetBytes("/Users/itin/go/src/dartech-course/lesson_2/wiki")
	for i := 0; i < len(slice); i++ {
		go CounterSentence(slice[i], ch)

		v, ok := <-ch
		if !ok {
			break
		} else {
			for a, b := range v {

				res[string(a)] += b
			}
		}
	}
	fmt.Printf("%v\n",res)
}


func CounterSentence(b []byte, ch chan map[byte]int) {
	mp := make(map[byte]int)
	for _, item := range b {
		if item >= 33 && item < 127 {
			if _, ok := mp[item]; ok {
				mp[item] += 1
			} else {
				mp[item] = 1
			}
		}
	}
	ch <- mp
}

func GetBytes(path string) [][]byte {
	b, err := ioutil.ReadFile("/Users/itin/go/src/dartech-course/lesson_2/wiki")
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	slice := bytes.Split(b, []byte("\n"))
	return slice
}