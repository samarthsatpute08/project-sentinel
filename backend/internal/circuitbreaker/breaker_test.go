package circuitbreaker

import (
	"sync"
	"testing"
	"time"
)

func TestBreakerConcurrency(t *testing.T) {
	b := New(Config{
		FailureThreshold: 5,
		Timeout:          time.Second,
	})

	var wg sync.WaitGroup
	// Simulate 100 concurrent failures
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.RecordFailure()
		}()
	}
	wg.Wait()

	if b.State() != StateOpen {
		t.Errorf("Expected state OPEN, got %v", b.State())
	}
}