package main

import (
	"fmt"
)

func readChan(ch chan int) {
	value := <-ch
	fmt.Println("CHAN VALUE: ", value)
}

func main() {
	fmt.Println("START MAIN")
	var ch chan int

	ch = make(chan int, 2)

	ch <- 100
	ch <- 200

	go readChan(ch)

	ch <- 300
	fmt.Println("END MAIN")
}
