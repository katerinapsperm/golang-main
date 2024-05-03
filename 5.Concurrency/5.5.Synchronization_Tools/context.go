package main

import (
	"context"
	"fmt"
	"time"
)

func sendData(ctx context.Context, num int) {
	timer := time.NewTimer(time.Duration(num) * time.Second)

	select {
	case <-ctx.Done():
		fmt.Printf("Процесс #%v отменен\n", num)
		return
	case <-timer.C:
		fmt.Printf("Данные процесса #%v успешно отправлены\n", num)
	}
}

func main() {
	parent := context.Background()
	ctx, _ := context.WithTimeout(parent, 3 * time.Second)

	for i := 1; i <= 10; i++ {
		go sendData(ctx, i)
	}
	time.Sleep(4 * time.Second)
}
