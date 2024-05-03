package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.Tick(1 * time.Second)

	count := 0
	for tick := range ticker {
		count++
		fmt.Printf("Tick #%v, time %v\n", count, tick)
	}
}
