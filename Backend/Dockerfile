# Start from Go base image
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o server ./cmd/server

# Small final image
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libc6
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]
