FROM golang:alpine AS builder

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go get -d -v
RUN go build

ENTRYPOINT ["./go_nhl"]