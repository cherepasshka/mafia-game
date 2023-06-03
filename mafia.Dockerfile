FROM golang:1.18-alpine

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /mafia-game/
COPY go.* /mafia-game/
RUN go mod download
ADD proto /mafia-game/proto/
ADD game-server /mafia-game/game-server/
ADD kafka-help /mafia-game/kafka-help/
RUN mkdir pdf
RUN go build -tags musl -o myapp  /mafia-game/game-server/cmd/main.go

