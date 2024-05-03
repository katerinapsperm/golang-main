package main

import (
	"fmt"
	"runtime"
	"time"
)

func write(ch chan<- int) {
	ch <- 1
}

func read(ch, quit <-chan int)  {
	for {
		select {
		case x := <-ch:
			fmt.Println("ch = ", x)
		case <-quit:
			fmt.Println("quit")
			return
		default:
			fmt.Println("default")
		}
	}
}

func main()  {
	ch := make(chan int)
	quit := make(chan int)

	go read(ch, quit)

	go write(ch)
	runtime.Gosched()

	go write(quit)
	time.Sleep(1 * time.Second)
}