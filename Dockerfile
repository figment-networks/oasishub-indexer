FROM golang:1.13 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o oasishub

#FROM golang:1.13.7-alpine3.10
#COPY --from=builder /app/oasishub /app/

EXPOSE 8081
#ENTRYPOINT ["/app/oasishub"]