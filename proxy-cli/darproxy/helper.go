package darproxy

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetPath(fd string) string {
	if fd == "" {
		fd = os.Getenv("HOME") + "/" + ".darproxy/config.json"
	}
	return fd
}

func GetResponseRoundRobin(count int, num int) ([]byte, error) {
	var body []byte
	var resp *http.Response
	var err error
	if count%2 == 0 {
		resp, err = http.Get(Cmd.config.Upstreams[num].Back[index0])
		if err != nil {
			return body, err
		}
	} else {
		resp, err = http.Get(Cmd.config.Upstreams[num].Back[index1])
		if err != nil {
			return body, err
		}
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func a(num int) ([]byte,error) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan []byte, 1)
	defer close(ch)
	go func(ctx context.Context, ch chan []byte) {
		resp, _ := http.Get(Cmd.config.Upstreams[num].Back[index0])
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		select {
		case ch <- body:
			cancel()
		case <-ctx.Done():
			fmt.Println("Canceled by timeout")
			return
		}
	}(ctx, ch)

	go func(ctx context.Context, ch chan []byte) {
		resp, _ := http.Get(Cmd.config.Upstreams[num].Back[index1])
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ch <- body
		select {
		case ch <- body:

			cancel()
		case <-ctx.Done():
			fmt.Println("Canceled by timeout")
			return
		}
	}(ctx, ch)
	select {
	case b := <- ch:
		return b, nil
	case <- time.After(5 * time.Second):
		return []byte("timeout 5 sec"), nil
	}
}

func GetResponseANYCAST(num int)([]byte,error) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan []byte, 1)

	f := func(num int, index int) {
		resp, _ := http.Get(Cmd.config.Upstreams[num].Back[index])
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		ch <- body

		select {
		case <-ctx.Done():
			fmt.Println("Canceled by timeout")
			return
		}
	}

	go f(num, index0)
	go f(num, index1)

	select {
	case b:=<-ch:
		return b,nil
	case <-time.After(1000 * time.Millisecond):
		cancel()
		return []byte("Timeout 5 sec"), nil
	}
}