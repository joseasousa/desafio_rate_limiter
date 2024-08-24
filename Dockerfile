FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -v -o server cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/server .
CMD ["./server"]
