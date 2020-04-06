FROM golang:1.13 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN GIT_SHA=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux

#Build server
RUN go build -a \
    -ldflags "-extldflags '-static' -w -s -X main.appSha=$GIT_SHA" \
    -o /go/bin/server \
    ./apps/server

#Build job
RUN go build -a \
    -ldflags "-extldflags '-static' -w -s -X main.appSha=$GIT_SHA" \
    -o /go/bin/job \
    ./apps/job

#Build cli
RUN go build -a \
    -ldflags "-extldflags '-static' -w -s -X main.appSha=$GIT_SHA" \
    -o /go/bin/cli \
    ./apps/cli

FROM alpine:3.10

COPY --from=builder /go/bin/server /go/bin/server
COPY --from=builder /go/bin/job /go/bin/job
COPY --from=builder /go/bin/cli /go/bin/cli

EXPOSE 8081
CMD ["/go/bin/server"]