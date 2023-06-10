FROM golang:1.18-alpine

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /mafia-game/
RUN mkdir pdf
RUN mkdir images
COPY go.* /mafia-game/
RUN go mod download
ADD scoreboard-service /mafia-game/scoreboard-service/
RUN go build -tags musl -o myapp  /mafia-game/scoreboard-service/cmd/main.go
