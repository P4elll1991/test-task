# test-task

Реализация тестового задания. Представляет собой 3 модуля:
 - http-server
 - broker
 - event-server

 Все три модуля развертываются в docker-compose. При этом event-server развернут в 2 сервиса с двумя разыми портами (9000 и 9001). Это необходимо для тестирования.

## Порядок тестирования:

- необходимо собрать модули : 
```bash
docker-compose build
```
- после удачной сборки их можно запустить : 
```bash
docker-compose up
```
- далее необходимо запустить тестовый http клиент : 
```bash
go run test-client/main.go
```
Он проверит каждый метод по 4 раза (сначала он будет проверять оба метода на удачные кейсы с валидными destination, а далее неудачные кейсы с destination, которые никуда не ведут). 

#### Пример успешного вывода теста:
```bash
$ go run test-client/main.go 
УСПЕШНЫЙ вызов sendChallenge с валидным destination
Ответ от сервера:
{"status":"OK"}
УСПЕШНЫЙ вызов sendMessage с валидным destination
Ответ от сервера:
{"status":"OK"}
УСПЕШНЫЙ вызов sendChallenge с не валидным destination
Ответ от сервера:
{"status":"OK"}
УСПЕШНЫЙ вызов sendMessage с не валидным destination
Ответ от сервера:
{"status":"OK"}
НЕ УСПЕШНЫЙ вызов sendChallenge с валидным destination
Ответ от сервера:
{"details":"rpc error: code = Unavailable desc = rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9002: connect: connection refused\"","error":"there is no active gRPC server at destination","status":"ERROR"}
НЕ УСПЕШНЫЙ вызов sendMessage с валидным destination
Ответ от сервера:
{"details":"rpc error: code = Unavailable desc = rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9002: connect: connection refused\"","error":"there is no active gRPC server at destination","status":"ERROR"}
НЕ УСПЕШНЫЙ вызов sendChallenge с не валидным destination
Ответ от сервера:
{"details":"rpc error: code = Unavailable desc = rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9003: connect: connection refused\"","error":"there is no active gRPC server at destination","status":"ERROR"}
НЕ УСПЕШНЫЙ вызов sendMessage с не валидным destination
Ответ от сервера:
{"details":"rpc error: code = Unavailable desc = rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9003: connect: connection refused\"","error":"there is no active gRPC server at destination","status":"ERROR"}
```
### http-server:
---
В HTTP сервере есть два хендлера:
- POST sendChallenge: в form-data подается переменная destination string, являющаяся ip:port серверного gRPC (event-server). Вызывает у broker по grpc процедуру sendChallenge и передает в нее destination
- POST sendMessage: в form-data подается destination string и data bytes. Вызывает у broker по gRPC процедуру sendMessage и и передает в нее destination и байты файла.

### Структура ответа:
---
JSON со 3 полями:
- status (OK | ERROR)
- error - текст ошибки
- details - детализация с трассировкой

Коды ответов:
- 200 - успех
- 400 - ошибка в параметрах (отсутствие активного gRPC сервера по адресу destination)
- 500 - любые другие фейлы

Необходимые переменные среды:
- BROKER_ADDR - адрес broker, куда будут перенаправляться запросы
- PORT - порт, на котором работает сервер (должен совпадать с аналогичным в broker)

### broker:
---
В broker есть 2 gRPC процедуры:
- sendChallenge: на вход принимает только переменную destination string, внутри генерирует рандомные 32 байта, пакует их в общий протобаф с названием Event и имеющим единственное поле bytes data и отправляет на серверный gRPC (event-server).
- sendMessage: на вход принимает destination.  Внутри генирируется произвольный RSA ключ, зашифроваются им данные data и запаковать в тот же Event.

Необходимые переменные среды:
- PORT - порт, на котором работает сервер

### event-server:
---
gRPC сервер имеет единственный метод:
- eventBus: на вход получает некий Event с переменной data, логгирует в консоль что был получен event с данными data.

Необходимые переменные среды:
- PORT - порт, на котором работает сервер