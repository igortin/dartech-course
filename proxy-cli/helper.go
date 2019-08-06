package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func GetPath(fd string) string {
	if fd == "" {
		fd = os.Getenv("HOME") + sep + ".darproxy/config.json"
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

func GetResponseANYCAST(num int) ([]byte, error) {
	ch0 := make(chan []byte)
	ch1 := make(chan []byte)
	go func(ch0 chan []byte) {
		resp, _ := http.Get(Cmd.config.Upstreams[num].Back[index0])
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		ch0 <- body
		defer close(ch0)
	}(ch0)
	go func(ch1 chan []byte) {
		resp, _ := http.Get(Cmd.config.Upstreams[num].Back[index1])
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		ch1 <- body
		defer close(ch1)
	}(ch1)

	select {
	case b := <-ch0:
		return b, nil
	case b := <-ch1:
		return b, nil
	}
}