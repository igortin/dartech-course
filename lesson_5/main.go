package main1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num, _ := strconv.Atoi(vars["id"])
	if num > 0 {
		a, b, c := getSlice(num + 1)
		fb := &FibNumber{a, b, c}
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
	r.HandleFunc("/fib/{id:[0-9]+}", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func fib(p int, m map[int]int) int {
	if p == 1 || p == 2 {
		m[p] = 1
		return 1
	}
	if p == 0 {
		m[p] = 0
		return 0
	}
	if w, t := m[p]; t {
		return w
	} else {
		m[p] = fib(p-1, m) + fib(p-2, m)
	}
	return m[p]
}

func getSlice(p int) (int, int, int) {
	m := make(map[int]int)
	fib(p+1, m)
	return m[p-2], m[p-1], m[p]
}