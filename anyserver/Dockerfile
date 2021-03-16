#FROM golang:latest
#
#COPY main.go /go/
#
#ENTRYPOINT go run /go/main.go

FROM golang:latest AS builder
WORKDIR /go/
#WORKDIR /go/src/github.com/alexellis/href-counter/
#RUN go get -d -v golang.org/x/net/html
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o anyserver .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/anyserver .
CMD ["./anyserver"]