FROM golang:1.16-alpine as builder

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

RUN go build -a -tags 'netgo osusergo' -o /go/bin/foaas-limiting main.go
RUN ldd /go/bin/foaas-limiting 2>&1 | grep -q 'Not a valid dynamic program'

LABEL description=foaas-limiting
LABEL builder=true
LABEL maintainer='Javier Garrone <javier3653@gmail.com>'

FROM alpine
COPY --from=builder go/bin/foaas-limiting /usr/local/bin

WORKDIR usr/local/bin
ENTRYPOINT [ "foaas-limiting", "serve" ]
EXPOSE 8080
