FROM golang:1.13.5 AS builder
WORKDIR /go/src/github.com/SlootSantos/janus-server
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 go build ./cmd/janus/main.go

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /go/src/github.com/SlootSantos/janus-server/main .
COPY --from=builder /go/src/github.com/SlootSantos/janus-server/.env .

EXPOSE 8888
CMD ["./main"] 