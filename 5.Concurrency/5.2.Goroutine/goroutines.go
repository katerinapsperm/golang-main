package main

import (
	"fmt"
	"runtime"
	"time"
)

func workAndPrint(num int) {
	fmt.Printf("START JOB #%v\n", num)
	var calc int
	for i:=0; i < 1000; i++ {
		calc = i * num
	}

	runtime.Gosched()
	fmt.Printf("END JOB #%v: calc = %v\n", num, calc)
}

func main()  {
	for i:=1; i <= 5; i++ {
		go workAndPrint(i)
	}

	time.Sleep(100 * time.Millisecond)
}
