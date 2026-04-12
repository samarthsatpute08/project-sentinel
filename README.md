# Project Sentinel

A high-performance API Multiplexer with a custom Circuit Breaker, 
real-time War Room dashboard, and chaos-tested infrastructure.

Built as part of the HCL GUVI Full-Stack Engineering Intern Assessment.

---

## What This Project Does

Imagine your app talks to an external API. That API sometimes goes down 
or gets slow. Without protection, every user request just hangs or fails.

Project Sentinel sits between your users and that API. It watches every 
request. The moment things go wrong, it trips a Circuit Breaker and 
silently reroutes all traffic to a backup API — without dropping a single 
request. A live dashboard shows you everything happening in real time.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (standard library only) |
| Frontend | React + Recharts |
| Containerisation | Docker + Docker Compose |
| Chaos Testing | Toxiproxy |
| Telemetry | WebSockets |

---

## How to Run It

Make sure you have Docker installed. Then run:

```bash
git clone https://github.com/samarthsatpute08/project-sentinel.git
cd project-sentinel
docker compose up --build
```

That's it. Docker spins up all 5 containers automatically:
- Go backend (circuit breaker + proxy router)
- React frontend (War Room dashboard)
- Primary API (dummy upstream)
- Fallback API (backup upstream)
- Toxiproxy (chaos injection layer)

Once running, open your browser:

| Service | URL | Notes |
|---------|-----|-------|
| War Room Dashboard | http://localhost:3000 | Open in browser |
| Backend Health Check | http://localhost:8080/health | Returns "ok" |
| Backend Proxy | http://localhost:8080/api/test | Hit with curl |
| Toxiproxy Control | http://localhost:8474 | REST API, use curl only |

---

## How to Trigger the Chaos Test

This proves the circuit breaker works under real hostile conditions.

**Step 1** — Inject 500ms latency and 20% packet loss into the Primary API:
```bash
./chaos.sh
```

**Step 2** — Send 100 requests through the router:
```bash
for i in $(seq 1 100); do curl -s http://localhost:8080/api/test; done
```

**Step 3** — Watch the dashboard at http://localhost:3000

You will see:
- Circuit Breaker badge flip from GREEN (CLOSED) to RED (OPEN)
- Traffic animation switch from Primary path to Fallback path
- Failure count increase in the stat cards

**Step 4** — Remove the chaos and watch the system recover:
```bash
curl -X DELETE http://localhost:8474/proxies/primary-api/toxics/latency-500ms
curl -X DELETE http://localhost:8474/proxies/primary-api/toxics/packet-loss-20pct
```
**Step 5** - Then send the traffic again
```bash
for i in $(seq 1 20); do curl -s http://localhost:8080/api/test; done
```
The circuit will move from OPEN → HALF-OPEN → CLOSED as it confirms 
the primary is healthy again.

---

## Architecture

```
Client
  │
  ▼
Go Router (port 8080)
  │  Circuit Breaker
  │  200ms Timeout
  │
  ├── NORMAL ──► Toxiproxy (8091) ──► Primary API (8081)
  │
  └── TRIPPED ──► Fallback API (8082)
  │
  ▼
React Dashboard (port 3000)
  ▲
  │ WebSocket
  └── Live Metrics (RPS, Latency, Circuit State)
```

All services run inside a single Docker Compose network.

---

## Key Design Decisions

### 1. Circuit Breaker (built from scratch)

The PDF required no external resilience libraries, so the circuit breaker 
is written entirely from scratch in Go.

It has three states:

- **CLOSED** — everything is normal, traffic goes to Primary API
- **OPEN** — too many failures detected, all traffic instantly goes to 
  Fallback API without even trying Primary
- **HALF-OPEN** — after a timeout, one probe request is allowed through 
  to test if Primary has recovered

State transitions use `sync/atomic` which means zero lock contention 
under high concurrency — multiple goroutines can read the state 
simultaneously without blocking each other.

### 2. 200ms Context Timeout

Every request to the Primary API is wrapped with Go's 
`context.WithTimeout(200ms)`. If the primary takes longer than 200ms 
to respond, the context is cancelled, a failure is recorded on the 
circuit breaker, and the request is immediately served by the Fallback 
API. The user never sees an error.

### 3. React Rendering Under High Frequency

The backend sends WebSocket updates dozens of times per second. If React 
re-rendered on every single message, the browser would freeze.

The solution: incoming WebSocket messages are stored in a `useRef` buffer 
(no re-render triggered). A `setInterval` running every 300ms drains the 
buffer and updates state — so React re-renders at most 3 times per second 
regardless of how fast the backend is sending data.

### 4. Memory Management (128MB limit)

The Go container runs with a hard 128MB memory limit. To stay within it:

- A single Hub goroutine manages all WebSocket clients instead of one 
  goroutine per client
- The broadcast channel is buffered (256 capacity) so the hub never 
  blocks the main request path
- Context cancellation on timeout prevents goroutine leaks when the 
  primary API is slow

---

## Project Structure

```
project-sentinel/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── circuitbreaker/         # Custom circuit breaker
│   │   ├── proxy/                  # HTTP reverse proxy router
│   │   └── telemetry/              # WebSocket hub
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/                     # Static assets
│   ├── src/
│   │   ├── assets/                 # Images and icons
│   │   ├── components/             # CircuitBadge, TrafficFlow, MetricsChart
│   │   ├── hooks/useMetrics.js     # WebSocket + buffered state
│   │   ├── App.jsx                 # War Room dashboard
│   │   └── main.jsx                # React entry point
│   ├── index.html
│   ├── package.json
│   └── Dockerfile
├── primary-api/                    # Dummy primary upstream
│   ├── main.go
│   └── Dockerfile
├── fallback-api/                   # Dummy fallback upstream
│   ├── main.go
│   └── Dockerfile
├── docker-compose.yml              # Full stack orchestration
├── chaos.sh                        # Toxiproxy chaos injection script
├── .gitignore
└── README.md
```
