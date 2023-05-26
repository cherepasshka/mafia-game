FROM golang:1.18-alpine

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /mafia-game/
COPY go.* /mafia-game/
RUN go mod download
ADD proto /mafia-game/proto/
ADD server /mafia-game/server/
ADD chat /mafia-game/chat/

RUN go build -tags musl -o myapp  /mafia-game/server/cmd/main.go
ENTRYPOINT ["/mafia-game/myapp"]
