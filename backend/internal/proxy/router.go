package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/samarthsatpute08/project-sentinel/internal/circuitbreaker"
	"github.com/samarthsatpute08/project-sentinel/internal/telemetry"
)

const primaryTimeout = 200 * time.Millisecond

type Router struct {
	primaryProxy  *httputil.ReverseProxy
	fallbackProxy *httputil.ReverseProxy
	breaker       *circuitbreaker.Breaker
	hub           *telemetry.Hub

	totalRequests atomic.Int64
	totalFailures atomic.Int64
	activeRoute   atomic.Value // stores string: "primary" or "fallback"

	// RPS tracking
	windowStart atomic.Int64
	windowCount atomic.Int64
}

func New(primaryURL, fallbackURL string, breaker *circuitbreaker.Breaker, hub *telemetry.Hub) (*Router, error) {
	pURL, err := url.Parse(primaryURL)
	if err != nil {
		return nil, err
	}
	fURL, err := url.Parse(fallbackURL)
	if err != nil {
		return nil, err
	}

	r := &Router{
		primaryProxy:  httputil.NewSingleHostReverseProxy(pURL),
		fallbackProxy: httputil.NewSingleHostReverseProxy(fURL),
		breaker:       breaker,
		hub:           hub,
	}
	r.activeRoute.Store("primary")
	r.windowStart.Store(time.Now().UnixMilli())

	// Override error handler so proxy errors are caught (not written as 502)
	r.primaryProxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("proxy: primary error: %v", err)
		r.breaker.RecordFailure()
		r.totalFailures.Add(1)
		r.activeRoute.Store("fallback")
		r.fallbackProxy.ServeHTTP(w, req)
	}

	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.totalRequests.Add(1)
	start := time.Now()

	if r.breaker.Allow() {
		// Primary path: wrap request with 200ms timeout
		ctx, cancel := context.WithTimeout(req.Context(), primaryTimeout)
		defer cancel()
		req = req.WithContext(ctx)

		r.activeRoute.Store("primary")
		r.primaryProxy.ServeHTTP(w, req)

		// If we got here without error handler firing, it was a success
		if ctx.Err() == nil {
			r.breaker.RecordSuccess()
		}
	} else {
		// Circuit is open — skip primary entirely
		r.activeRoute.Store("fallback")
		r.fallbackProxy.ServeHTTP(w, req)
	}

	latency := float64(time.Since(start).Milliseconds())
	r.publishMetrics(latency)
}

func (r *Router) publishMetrics(latency float64) {
	now := time.Now().UnixMilli()
	windowStart := r.windowStart.Load()
	elapsed := float64(now-windowStart) / 1000.0

	var rps float64
	count := r.windowCount.Add(1)
	if elapsed > 0 {
		rps = float64(count) / elapsed
	}

	// Reset window every 5 seconds
	if elapsed >= 5.0 {
		r.windowStart.Store(now)
		r.windowCount.Store(0)
	}

	r.hub.Publish(telemetry.Metrics{
		RPS:           rps,
		ActiveRoute:   r.activeRoute.Load().(string),
		CircuitState:  r.breaker.State().String(),
		TotalRequests: r.totalRequests.Load(),
		TotalFailures: r.totalFailures.Load(),
		Latency:       latency,
	})
}