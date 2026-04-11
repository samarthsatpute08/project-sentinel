FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init fallback && CGO_ENABLED=0 go build -o api main.go

FROM alpine:latest
RUN apk add --no-cache wget
COPY --from=builder /app/api /api
EXPOSE 8082
ENTRYPOINT ["/api"]