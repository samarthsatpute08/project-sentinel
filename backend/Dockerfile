FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o sentinel \
    ./cmd/server

FROM alpine:latest
RUN apk add --no-cache wget

COPY --from=builder /app/sentinel /sentinel

EXPOSE 8080

ENTRYPOINT ["/sentinel"]