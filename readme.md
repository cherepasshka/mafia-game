# Service oriented architecture practice-2

- [Описание](#description)
- [Использование](#usage)
- [Профили пользователей и взаимодействие с http](#http)

<a name="description"></a> 
## Сетевая версия игры Мафия

Сервис для игры в мафию состоит из сервера, расположенного в [server](https://github.com/cherepasshka/mafia-game/tree/main/server), и клиента, расположенного в [client](https://github.com/cherepasshka/mafia-game/tree/main/client). Скомпилированный файл клиента, готовый к запуску, это файл [mafia-client](https://github.com/cherepasshka/mafia-game/blob/main/mafia-client). Докер образ сервера выложен в [dockerhub](https://hub.docker.com/repository/docker/cherepashka/soa-practice-2)

<a name="usage"></a> 
## Использование

### 0. Генерация `proto` шаблона
```bash
protoc --go_out=. proto/mafia-game.proto --go-grpc_out=.
```

### 1. Запуск сервера
По умолчанию сервер поднимается на 9000 порту TCP

Запуск с помощью `docker-compose`:
```bash
docker-compose up
```
Кафка довольно долго поднимается (около 2-3 минут), поэтому нужно подождать перед тем, как запускать клиента.

### 2. Запуск клиента
Можно указать флаг для порта `--port 9000` и для хоста `--host 127.0.0.1`
```bash
go run client/cmd/main.go
```

Для сессии игры в мафию необходимо 4 игрока: 2 мирных жителя, комиссар и мафия.

<a name="http"></a> 
## Профили пользователей и взаимодействие с http
По умолчанию на порту 9001 (можно поменять через переменную окружения в docker-compose.yaml) поднимается http сервер.

### Добавление пользователя
Добавить пользователя можно с помощью `curl`, важно, что поля пользователя (логин, гендер, почту, картинку) нужно передавать как переменную формы (флаг `-F`).
Создание пользователя с логином `user_1`, почтой `user_1@example.com`, гендером `female` и с аватаром, который находится в файле на компьютере клиента `picture.jpg` (поддерживаются только картинки формата jpg):
```bash
curl -X POST -F "image=@picture.jpg" -F "login=user_1" -F "email=user_1@example.com" -F "gender=female" http://localhost:9001/users/1
```
Любое из полей можно опустить.
### Обновление пользователя
Выполняется аналогично, но используется http метод `PUT`
```bash
curl -X PUT -F "image=@picture.jpg" -F "login=user_1" -F "email=user_1@example.com" -F "gender=female" http://localhost:9001/users/1
```
### Удаление пользователя
Удаление пользователя с логином `user_1`:
```bash
curl -X DELETE http://127.0.0.1:9001/users/1
```
### Получение пользователя
Получение ссылки на pdf страницу профиля пользователя с логином `user_1`:
```bash
curl -X GET http://127.0.0.1:9001/users/1
```
Получение ссылки на pdf страницу профилей пользователей с логином `user_1`, `user_2` и `user_3` (таким образом можно получить профили сколки угодно пользователей):
```bash
curl -X GET http://127.0.0.1:9001/users/?logins=user_1,user_2,user_3
```
