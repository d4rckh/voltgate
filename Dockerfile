# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o voltgate-proxy ./main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/voltgate-proxy /app/voltgate-proxy

# Run the proxy
CMD ["/app/voltgate-proxy"]
