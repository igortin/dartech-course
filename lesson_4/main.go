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
	inputFile = "/tmp/bigfile"
	outputFile = "/tmp/outfilename"
)



type FdWriter struct{
	fd io.Writer
}

func (wr FdWriter) Write(data []byte) (n int, err error) {
	tmp := append([]byte(time.Now().String()),data...)
	tmp = append(tmp, []byte("\n")...)
	return wr.fd.Write(tmp)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan []byte)
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	defer close(ch)
	defer close(ch1)
	defer close(ch2)

	inputFilename, err := os.OpenFile(inputFile,  os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer inputFilename.Close()


	outputFilename, err := os.OpenFile(outputFile,  os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFilename.Close()
	

	stdOut := &FdWriter{
		os.Stdout,
		}

	fileOut := &FdWriter{
		outputFilename,
	}

	go func(ch chan<- []byte) {
		wg.Add(1)
		line := make([]byte, 256)
		reader := bufio.NewReader(inputFilename)
		for {
			_, err = reader.Read(line)
			if err != nil || err == io.EOF {
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
				bw := bufio.NewWriterSize(stdOut, 256) // принимает объект который реализует интерфейс io.Writer
				_, _ = bw.Write(val)						 // вызов метода Write объекта *File который находится под stdOut
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
			time.Sleep(1000 * time.Millisecond)
			if val, opened := <-ch2; opened {
				bw := bufio.NewWriterSize(fileOut, 256)  // принимает объект который реализует интерфейс io.Writer
				_, _ = bw.Write(val)						   // вызов метода Write объекта *File который находится под fileOut
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