FROM golang:alpine AS builder

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN cd cmd && go get -d -v
RUN go build -o go_nhl ./cmd/main.go

ENTRYPOINT ["./go_nhl"]