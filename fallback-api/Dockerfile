FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init fallback && CGO_ENABLED=0 go build -o api main.go

FROM scratch
COPY --from=builder /app/api /api
EXPOSE 8082
ENTRYPOINT ["/api"]