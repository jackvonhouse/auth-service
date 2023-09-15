# Сервис аутентификации

## Запуск

```
go run ./cmd/main.go [-config <путь>]
```

## HTTP API

| Метод | Эндпоинт | Дополнительно
| ----------- | ----------- | ----------- |
| POST | /api/v1/auth | ?guid={guid} |
| POST | /api/v1/refresh | {"access_token": "...", "refresh_token": "..."} |

### /api/v1/auth

Пример запроса

```
curl -X POST -i '127.0.0.1:8081/api/v1/auth?guid=773697bb-3c65-459c-8aaa-d3cb5e90233g'
```

Пример ответа

``` json
{
	"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJndWlkIjoiNzczNjk3YmItM2M2NS00NTljLThhYWEtZDNjYjVlOTAyMzNnIiwiaWQiOiI2NTA0MzJmNzhmZmNmZmI0MDk4NDI1NWMiLCJleHAiOjE2OTQ3Nzc2MDd9.3kmP-49loat_CFN4afYQeAt1tMzcHzuw_6GMW56lEk6NHzjY_vuCSn9hp0EQ_jPhD8bapHpx-bxKAgl46Jkz9Q",
	"refresh_token":"FuQtQ3AJ6Zqe1j4YHQMXAxFLmQ2R1ISr7Zm8O53cDKY"
}
```

### /api/v1/refresh

Пример запроса

```
curl -X POST -i '127.0.0.1:8081/api/v1/refresh' --data '{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJndWlkIjoiNzczNjk3YmItM2M2NS00NTljLThhYWEtZDNjYjVlOTAyMzNnIiwiaWQiOiI2NTA0MzJmNzhmZmNmZmI0MDk4NDI1NWMiLCJleHAiOjE2OTQ3Nzc2MDd9.3kmP-49loat_CFN4afYQeAt1tMzcHzuw_6GMW56lEk6NHzjY_vuCSn9hp0EQ_jPhD8bapHpx-bxKAgl46Jkz9Q","refresh_token":"FuQtQ3AJ6Zqe1j4YHQMXAxFLmQ2R1ISr7Zm8O53cDKY"}'
```

Пример ответа

``` json
{
	"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJndWlkIjoiNzczNjk3YmItM2M2NS00NTljLThhYWEtZDNjYjVlOTAyMzNnIiwiaWQiOiI2NTA0MzJmNzhmZmNmZmI0MDk4NDI1NWMiLCJleHAiOjE2OTQ3Nzc2MDd9.3kmP-49loat_CFN4afYQeAt1tMzcHzuw_6GMW56lEk6NHzjY_vuCSn9hp0EQ_jPhD8bapHpx-bxKAgl46Jkz9Q",
	"refresh_token":"FuQtQ3AJ6Zqe1j4YHQMXAxFLmQ2R1ISr7Zm8O53cDKY"
}
```

## База данных

В качестве инстанса использовалась облачная база данных MongoDB - [**Atlas**](https://www.mongodb.com/atlas).

Для корректной работы необходимо создать базу данных ```auth``` с коллекцией ```refresh_tokens```.

## JWT-токены

| Тип токена | Алгоритм хеширования | Хранение в БД |
| ----------- | ----------- | ----------- |
| Access | SHA512 | Нет
| Refresh | base64(<случайные байты>)[^1] | Да

[^1]: Перед добавлением токена в базу данных он хешируется при помощи пакета ```crypto/bcrypt```.

**Refresh** токен генерируется при помощи пакета ```crypto/rand``` посредством чтения из потока N-го количества случайных байт.

В отличие от **Refresh**, **Access** токен хранит в себе ```GUID``` и ```id```.

### Обратное связывание токенов

При генерации пары токенов **Refresh** токен добавляетсяв базу данных, а полученный ```objectID``` добавляется в **Access** токен под аргументом ```id```.

При проверке токенов сервис возьмёт первый **Refresh** токен (переданнный пользователем), получит второй токен (из базы данных по ```id``` из **Access** токена), захеширует второй и сверит оба токена.

### Срок истечения токенов

Если для **Access** токена срок истечения можно указать внутри его структуры, то в случае с **Refresh** токеном в базе данных создаётся TTL-индекс с полем ```expired_at```, который инициирует удаление токена при истечении срока.