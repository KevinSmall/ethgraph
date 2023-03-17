package work

import (
	"sync/atomic"
	"testing"
	"time"
)

type testWorker struct {
	executed int32
}

func (t *testWorker) Task() {
	atomic.AddInt32(&t.executed, 1)
}

func TestPool_Run(t *testing.T) {
	maxGoroutines := 5
	pool := New(maxGoroutines)
	tw := &testWorker{}

	// Submit test workers to the pool
	for i := 0; i < 10; i++ {
		pool.Run(tw)
	}

	// Give some time for the tasks to be executed
	time.Sleep(1 * time.Second)

	// Check if all tasks have been executed
	if tw.executed != 10 {
		t.Errorf("Expected 10 tasks to be executed, but got %d", tw.executed)
	}

	pool.Shutdown()
}
