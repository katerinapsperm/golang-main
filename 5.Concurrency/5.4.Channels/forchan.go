package main

import "fmt"

func writeChan(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	//close(ch)
}

func main() {
	fmt.Println("START MAIN")
	ch := make(chan int)

	go writeChan(ch)

	for i := range ch {
		fmt.Println("chan i = ", i)
	}
}
