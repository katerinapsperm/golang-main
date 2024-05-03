package main

import (
	"fmt"
	"time"
)

func jobWithTimeout(t *time.Timer, q chan int) {
	time.Sleep(1 * time.Second)
	select {
	case <-t.C:
		fmt.Println("Время вышло")
	case <-q:
		if !t.Stop() {
			<-t.C
		}
		fmt.Println("Таймер остановлен")
	default:
		fmt.Println("Завершение функции")
	}
}

func main() {
	timer := time.NewTimer(1 * time.Second)
	quit := make(chan int)

	go jobWithTimeout(timer, quit)

	time.Sleep(2 * time.Second)
}
