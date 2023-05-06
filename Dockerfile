FROM golang:1.18-alpine

WORKDIR /mafia-game/
COPY go.* /mafia-game/
RUN go mod download
ADD server /mafia-game/server/
ADD proto /mafia-game/proto/