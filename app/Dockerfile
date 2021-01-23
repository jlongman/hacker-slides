FROM golang:1.12

ENV VERSION=1.0.0

LABEL version=$VERSION

WORKDIR $GOPATH/src/github.com/msoedov/hacker-slides

ENV GO111MODULE on
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLE=0 go build  -a -tags netgo -ldflags '-w -extldflags "-static"' -o /bin/app *.go
