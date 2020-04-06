FROM golang:1.13 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY ./apps/server .
RUN go build -o /go/bin/server

COPY ./apps/job .
RUN go build -o /go/bin/job

COPY ./apps/cli .
RUN go build -o /go/bin/cli

FROM alpine:3.10

COPY --from=builder /go/bin/server /go/bin/server
COPY --from=builder /go/bin/job /go/bin/job
COPY --from=builder /go/bin/cli /go/bin/cli

EXPOSE 8081
ENTRYPOINT ["/go/bin/server"]