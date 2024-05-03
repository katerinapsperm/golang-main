package main

import (
	"fmt"
	"slurm/client"
	"sync/atomic"
	"time"
)

func addAge(cl *client.Client, add int) {
	atomic.AddInt64(&cl.Age, int64(add))
}

func main()  {
	cl := &client.Client{}

	for i := 1; i <= 1000; i++ {
		go addAge(cl, 1)
	}

	time.Sleep(1*time.Second)
	fmt.Printf("%#v\n", cl)
}
