# Service oriented architecture practice-2

## Сетевая версия игры Мафия

Сервис для игры в мафию состоит из сервера, расположенного в [server](https://github.com/cherepasshka/mafia-game/tree/main/server), и клиента, расположенного в [client](https://github.com/cherepasshka/mafia-game/tree/main/client). Скомпилированный файл клиента, готовый к запуску, это файл [mafia-client](https://github.com/cherepasshka/mafia-game/blob/main/mafia-client). Докер образ сервера выложен в [dockerhub](https://hub.docker.com/repository/docker/cherepashka/soa-practice-2)

## Использование

### 0. Генерация `proto` шаблона
```bash
protoc --go_out=. proto/mafia-game.proto --go-grpc_out=.
```

### 1. Запуск сервера
По умлочанию сервер поднимается на 9000 порту TCP

С помощью `docker-compose`:
```bash
docker-compose up
```
Запуск без контейнера:
```bash
go run server/cmd/main.go
```
### 2. Запуск клиента
Можно указать флаг для порта `--port 9000` и для хоста `--host 127.0.0.1`
```bash
./mafia-client
```

Для сессии игры в мафию необходимо 4 игрока: 2 мирных жителя, комиссар и мафия.