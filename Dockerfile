FROM golang:1.16.4-alpine3.13

RUN mkdir -p /go/src/github.com/rtgnx/gotty
WORKDIR /go/src/github.com/rtgnx/gotty
COPY . .
RUN go mod tidy && go mod vendor
RUN go build -o /usr/bin/gotty
RUN apk add docker
RUN chmod +x shell.sh && cp shell.sh /usr/bin/shell
WORKDIR /
CMD ["/usr/bin/gotty", "--port", "$PORT"]