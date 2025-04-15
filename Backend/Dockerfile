# Start from Go base image
FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o server ./cmd/server

# Small final image
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]
