package darproxy

import (
	"io/ioutil"
	"net/http"
	"time"
	"context"
)

var (
	count = &Count{0}
)

func GetResponseRoundRobin(index int, config Config) ([]byte, error) {
	var body []byte
	var resp *http.Response
	var err error

	if count.scorer == len(config.Upstreams) {
		count.ResetCount()
	}

	upstream := count.scorer % len(config.Upstreams)

	resp, err = http.Get(config.Upstreams[index].Back[upstream]) // NO
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	count.IncreaseCount()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func GetResponseAnycast(index int, config Config)([]byte,error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []byte, 1)
	f := func(index int, n int) {
		resp, _ := http.Get(config.Upstreams[index].Back[n])
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		ch <- body
		select {
		case <-ctx.Done():
			return
		}
	}
	go f(index, 0)
	go f(index, 1)
	select {
	case b:=<-ch:
		return b,nil
	case <-time.After(1000 * time.Millisecond):
		return []byte("Timeout 5 sec"), nil
	}
}