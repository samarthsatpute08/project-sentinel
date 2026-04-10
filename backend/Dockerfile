# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -ldflags="-s -w" strips debug info → smaller binary
# CGO_ENABLED=0 → fully static binary, works in scratch
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o sentinel \
    ./cmd/server

# Stage 2: Runtime (scratch = no OS, just our binary)
FROM scratch

COPY --from=builder /app/sentinel /sentinel

EXPOSE 8080

ENTRYPOINT ["/sentinel"]