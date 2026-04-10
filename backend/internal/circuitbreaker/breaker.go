package circuitbreaker

import (
	"sync"
	"sync/atomic"
	"time"
)

type State int32

const (
	StateClosed   State = iota // Normal: requests go to primary
	StateOpen                  // Tripped: requests go to fallback immediately
	StateHalfOpen              // Testing: one probe request allowed through
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

type Config struct {
	FailureThreshold int           // How many failures before tripping
	SuccessThreshold int           // How many successes in half-open before closing
	Timeout          time.Duration // How long to stay Open before trying half-open
}

type Breaker struct {
	cfg            Config
	state          atomic.Int32
	failures       atomic.Int32
	successes      atomic.Int32
	lastFailureAt  time.Time
	mu             sync.Mutex
	OnStateChange  func(from, to State)
}

func New(cfg Config) *Breaker {
	return &Breaker{cfg: cfg}
}

// Allow returns true if the request should go to the primary API
func (b *Breaker) Allow() bool {
	state := State(b.state.Load())

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		b.mu.Lock()
		elapsed := time.Since(b.lastFailureAt)
		b.mu.Unlock()
		if elapsed >= b.cfg.Timeout {
			// Try transitioning to half-open
			if b.state.CompareAndSwap(int32(StateOpen), int32(StateHalfOpen)) {
				b.notify(StateOpen, StateHalfOpen)
			}
			return true // Allow one probe
		}
		return false
	case StateHalfOpen:
		return true // Let the probe through
	}
	return false
}

// RecordSuccess records a successful request
func (b *Breaker) RecordSuccess() {
	state := State(b.state.Load())
	if state == StateHalfOpen {
		n := b.successes.Add(1)
		if int(n) >= b.cfg.SuccessThreshold {
			b.reset()
		}
	}
}

// RecordFailure records a failed or timed-out request
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	b.lastFailureAt = time.Now()
	b.mu.Unlock()

	n := b.failures.Add(1)
	state := State(b.state.Load())

	if state == StateHalfOpen {
		// Any failure in half-open snaps back to open immediately
		b.trip()
		return
	}

	if state == StateClosed && int(n) >= b.cfg.FailureThreshold {
		b.trip()
	}
}

func (b *Breaker) trip() {
	if b.state.CompareAndSwap(int32(StateClosed), int32(StateOpen)) ||
		b.state.CompareAndSwap(int32(StateHalfOpen), int32(StateOpen)) {
		b.successes.Store(0)
		b.notify(StateClosed, StateOpen)
	}
}

func (b *Breaker) reset() {
	b.state.Store(int32(StateClosed))
	b.failures.Store(0)
	b.successes.Store(0)
	b.notify(StateHalfOpen, StateClosed)
}

func (b *Breaker) State() State {
	return State(b.state.Load())
}

func (b *Breaker) Failures() int {
	return int(b.failures.Load())
}

func (b *Breaker) notify(from, to State) {
	if b.OnStateChange != nil {
		go b.OnStateChange(from, to)
	}
}