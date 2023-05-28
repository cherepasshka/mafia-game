FROM golang:1.18-alpine

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /mafia-game/
COPY go.* /mafia-game/
RUN go mod download
ADD proto /mafia-game/proto/
ADD chat-server /mafia-game/chat-server/
ADD kafka-help /mafia-game/kafka-help/

RUN go build -tags musl -o myapp  /mafia-game/chat-server/cmd/main.go

