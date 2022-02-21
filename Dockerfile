FROM golang:1.15-alpine as builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://proxy.golang.org,direct"

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download -x
COPY . .
COPY *.csv /go/bin/

RUN go build -a -tags 'netgo osusergo' -o /go/bin/user-retention main.go

LABEL description=user-retention
LABEL builder=true
LABEL maintainer='Santiago Leira'

FROM alpine
COPY --from=builder go/bin/user-retention /usr/local/bin
COPY --from=builder go/bin/*.csv /usr/local/bin/

WORKDIR usr/local/bin
ENTRYPOINT [ "user-retention", "script" ]
EXPOSE 8080
