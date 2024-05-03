package main

import (
	"fmt"
	"sync"
	"time"
)

func sleep(t time.Duration, wg *sync.WaitGroup) {
	fmt.Println("Sleep ", t)
	time.Sleep(t)
	wg.Done()
}

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go sleep(1*time.Second, wg)

	wg.Add(1)
	go sleep(2*time.Second, wg)

	wg.Add(1)
	go sleep(3*time.Second, wg)

	wg.Add(1)
	wg.Wait()
	fmt.Println("End main")
}
