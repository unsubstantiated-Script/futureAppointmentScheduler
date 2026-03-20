FROM golang:1.25.8 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY migrations ./migrations
COPY data ./data

EXPOSE 8080

CMD ["./server"]