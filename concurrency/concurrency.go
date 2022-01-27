package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	HandleTimeout        = 2
	HandleTaskTimeout    = 3
	HandleRequestTimeout = 60
)

type Response struct {
	Task   string `json:"task"`
	Status string `json:"status"`
}

func ConcurrentProcessingOfTasks(ctx context.Context,
	tasks []string) []Response {
	// Create a deadline to wait for.
	ctxWithTimeout, cancel := context.WithTimeout(ctx,
		time.Duration(HandleTimeout)*time.Second)

	var wg sync.WaitGroup
	taskProcessInfo := make(chan Response, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go HandleTask(ctxWithTimeout, task, taskProcessInfo, &wg)
	}

	done := make(chan bool)

	/*
		This goroutine helps in figuring out if all goroutines above
		ran successfully or not.
	*/
	go func() {
		wg.Wait()
		close(taskProcessInfo)
		done <- true
	}()

	/*
		This is used to cancel the ctx if either the wg marks all
		goroutines as done or if request timeout happens  which
		will lead to us returning a nil response and the goroutines
		will run in background till HandleTimeout (timeout of the
		the ctx sent to the goroutines).
	*/
	select {
	case <-done:
		fmt.Println("req successful")
		close(done)
		cancel()
	case <-time.After(time.Duration(HandleRequestTimeout) * time.Millisecond):
		fmt.Println("req timeout")
		return nil
	}

	length := len(taskProcessInfo)

	response := make([]Response, 0, length)
	for info := range taskProcessInfo {
		response = append(response, info)
	}

	return response
}

func HandleTask(ctx context.Context,
	task string,
	taskProcessInfo chan Response,
	wg *sync.WaitGroup) {

	defer func() {
		fmt.Println("task done: ", task)
		wg.Done()
	}()

	var response = Response{
		Task:   task,
		Status: "success",
	}

	select {
	case <-ctx.Done():
		fmt.Println("ctx is finished for task: ", task)
		return
	case <-time.After(time.Duration(HandleTaskTimeout) * time.Second):
		taskProcessInfo <- response
	}
}
