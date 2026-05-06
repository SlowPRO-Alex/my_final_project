FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o todo-app main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/todo-app ./todo-app
COPY web ./web
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD="123456"
VOLUME /data
CMD ["./todo-app"]