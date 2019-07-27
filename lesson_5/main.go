package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	num, err := strconv.Atoi(vars["id"])

	if err != nil {
		panic("request has no parameter id")
	}

	if num > 0 {

		slice := []int{}

		ch := make(chan int)
		defer close(ch)

		go func(ch chan<- int) {
			ch <- fib(num - 1)
		}(ch)

		go func(ch chan<- int) {
			ch <- fib(num)
		}(ch)

		go func(ch chan<- int) {
			ch <- fib(num + 1)
		}(ch)

		for i := 0; i < 3; i++ {
			slice = append(slice, <-ch)
		}

		sort.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})

		fb := &FibNumber{slice[0], slice[1], slice[2]}

		fbj, _ := json.Marshal(fb)

		_, _ = w.Write(fbj)
	} else {
		_, _ = w.Write([]byte("Wrong number! please try again!"))
	}
}

type FibNumber struct {
	Prev    int
	Current int
	Next    int
}




func main() {
	r := mux.NewRouter()
	r.HandleFunc("/fib/{id}", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func fib(p int) int {
	if p == 1 || p == 2 {
		return 1
	}
	if p == 0 {
		return 0
	}
	return fib(p-1) + fib(p-2)
}