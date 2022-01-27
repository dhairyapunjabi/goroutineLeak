package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dhairyapunjabi/goroutineLeak/concurrency"
)

func main() {
	tasks := []string{"task1", "task2"}
	fmt.Println("req")

	concurrency.ConcurrentProcessingOfTasks(context.Background(), tasks)

	time.Sleep(6 * time.Second)
}
