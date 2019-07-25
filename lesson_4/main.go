package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	path = "/tmp/bigfile"
	pathWrite = "/tmp/WriteTo"
)

type ConsoleWriter struct{}


func (wr ConsoleWriter) Write(data []byte) (n int, err error) {
	fmt.Printf("%v %v\n", time.Now(), string(data))
	return len(data), nil
}

type FdWriter struct{}

func (wr FdWriter) Write(data []byte) (n int, err error) {
	file, err := os.OpenFile(pathWrite, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, _ = file.WriteString(string(append([]byte(time.Now().String()),data...)))
	return len(data), nil
}

func main() {
	var mnOnce sync.Once
	var wg sync.WaitGroup
	ch := make(chan []byte)
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)
	defer close(ch)
	defer close(ch1)
	defer close(ch2)

	mnOnce.Do(func() {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			_, err := os.Create(path)
			if err != nil {
				os.Exit(1)
			}
		}
		_, er := os.Stat(pathWrite)
		if os.IsNotExist(er) {
			_, er := os.Create(pathWrite)
			if er != nil {
				os.Exit(1)
			}
		}
	})

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	go func(ch chan<- []byte) {
		wg.Add(1)
		line := make([]byte, 256*1024)
		reader := bufio.NewReader(file)
		for {
			_, err = reader.Read(line)
			if err == io.EOF {
				break
			}
			ch <- line

		}
		wg.Done()
	}(ch)

	go func(ch <-chan []byte, ch1 chan<- []byte,  ch2 chan<- []byte) {
		wg.Add(1)
		for {
			if val, opened := <-ch; opened {
				ch1 <- val
				ch2 <- val
			} else {
				break
			}
		}
		wg.Done()
	}(ch, ch1, ch2)

	go func(ch1 <-chan []byte) {
		wg.Add(1)
		for {
			time.Sleep(1000 * time.Millisecond)
			if val, opened := <-ch1; opened {
				bw := bufio.NewWriterSize(&ConsoleWriter{}, 256*1024)
				_, _ = bw.Write(val)
				_ = bw.Flush()
			} else {
				break
			}
		}
		wg.Done()
	}(ch1)

	go func(ch2 <-chan []byte) {
		wg.Add(1)
		for {
			if val, opened := <-ch2; opened {
				time.Sleep(1000 * time.Millisecond)
				bw := bufio.NewWriterSize(&FdWriter{}, 256 * 1024)
				_, _ = bw.Write(val)
				_ = bw.Flush()
			} else {
				break
			}
		}
		wg.Done()
	}(ch2)
	fmt.Printf("")
	wg.Wait()
}