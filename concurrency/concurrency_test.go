package concurrency_test

import (
	"context"
	"testing"

	"github.com/fortytw2/leaktest"

	"github.com/dhairyapunjabi/goroutineLeak/concurrency"
)

func TestConcurrentProcessingOfTasks(t *testing.T) {

	defer leaktest.Check(t)()

	tasks := []string{"task1", "task2"}

	concurrency.ConcurrentProcessingOfTasks(context.Background(), tasks)
}
