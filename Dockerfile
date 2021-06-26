FROM golang:1.16.4-alpine3.13

RUN mkdir -p /go/src/github.com/rtgnx/gotty
WORKDIR /go/src/github.com/gotty/gotty
COPY . .
RUN go mod tidy && go mod vendor
RUN go build -o /usr/bin/gotty
CMD ["/usr/bin/gotty"]