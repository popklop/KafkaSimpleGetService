# Order Service 

**Order Service** — сервис для получения, хранения и предоставления заказов.

Он читает заказы из **Kafka**, сохраняет их в **PostgreSQL**, кэширует в памяти и предоставляет **HTTP API** и простой **веб‑интерфейс** для доступа к данным. Также сервис экспортирует метрики для **Prometheus** и **Grafana**.

---

## Возможности

* Приём заказов из Kafka
* Хранение данных в PostgreSQL
* In‑memory кэш для быстрого доступа
* HTTP API для получения заказов
* Веб‑интерфейс для поиска заказа по ID
* Метрики Prometheus
* Готовая инфраструктура через Docker Compose

---

## Стек технологий

* **Go** 1.25
* **PostgreSQL** 15
* **Apache Kafka**
* **Prometheus**
* **Grafana**
* **Docker / Docker Compose**
* **Cassowary** — нагрузочное тестирование

---

## Архитектура

Проект реализован по принципам **Clean Architecture**.

### Слои

* **domain** — бизнес‑сущности и правила
* **usecase** — сценарии использования
* **infrastructure** — внешние адаптеры (БД, Kafka, HTTP, кэш)
* **app** — инициализация и запуск приложения
* **config** — конфигурация из переменных окружения
* **metrics** — метрики Prometheus

### Структура проекта

```text
order-service/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── domain/
│   ├── infrastructure/
│   │   ├── cache/
│   │   ├── http/
│   │   ├── kafka/
│   │   ├── postgres/
│   │   └── web/
│   └── usecase/
├── metrics/
├── scripts/
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

---

## Запуск проекта

### 1. Запуск инфраструктуры

```bash
docker-compose up -d
```

Контейнеры:

| Контейнер  | Порт                 |
| ---------- | -------------------- |
| PostgreSQL | 5432                 |
| Zookeeper  | 2181                 |
| Kafka      | 9092                 |
| Prometheus | 9090                 |
| Grafana    | 3000 (admin / admin) |

---

### 2. Запуск сервиса

```bash
go run cmd/app/main.go
```

Сервис будет доступен по адресу:

```
http://localhost:8081
```

---

## Конфигурация

Настраивается через переменные окружения:

| Переменная   | Значение по умолчанию | Описание        |
| ------------ | --------------------- | --------------- |
| DB_USER      | postgres              | Пользователь БД |
| DB_PASSWORD  | order_pass            | Пароль БД       |
| DB_HOST      | localhost             | Хост БД         |
| DB_PORT      | 5432                  | Порт БД         |
| DB_NAME      | wborders              | Имя БД          |
| KAFKA_BROKER | localhost:9092        | Kafka брокер    |
| KAFKA_TOPIC  | orders                | Kafka топик     |
| KAFKA_GROUP  | order-service         | Consumer group  |
| HTTP_PORT    | 8081                  | HTTP порт       |

При запуске через `docker-compose` значения уже настроены.

---

## Отправка тестовых заказов в Kafka

Используется скрипт `producer/main.go`, который генерирует случайные заказы через fakeit/v7.

```bash
cd producer
go run main.go -n 5
```

Параметры:

* `-n` — количество сообщений
* `-broker` — адрес Kafka
* `-topic` — Kafka топик

---

## HTTP API

### Получить заказ по ID

```http
GET /order/{order_uid}
```

Пример:

```bash
curl http://localhost:8081/order/b10198ce-653f-4c04-b038-360fb4f43f55
```

Ответ:

```json
{
  "ID": "b10198ce-653f-4c04-b038-360fb4f43f55",
  "TrackNumber": "TRACK-qFsbBl",
  "Entry": "WBIL",
  "Delivery": {
    "Name": "Brady Walters",
    "Phone": "7350512683",
    "Zip": "31072",
    "City": "Memphis",
    "Address": "102 New Extensionbury, 85185",
    "Region": "California",
    "Email": "albertogriffin@butler.info"
  },
  "Payment": {
    "Transaction": "dc5c2814-8049-40f5-adaf-909349624077",
    "RequestID": "25078978-63da-45ac-a529-4feee6ff9052",
    "Currency": "GBP",
    "Provider": "wbpay",
    "Amount": 4821,
    "PaymentAt": "2026-03-03T08:47:40Z",
    "Bank": "Zonability",
    "DeliveryCost": 469,
    "GoodsTotal": 4427,
    "CustomFee": 0
  },
  "Items": [
    {
      "ChrtID": 9699,
      "TrackNumber": "TRACK-qFsbBl",
      "Price": 131,
      "RID": "94f6e2d1-15cf-4a01-aa94-3dc92989acf0",
      "Name": "Dumbbell Mind",
      "Sale": 2,
      "Size": "L",
      "NmID": 12740,
      "Brand": "Archimedes Inc.",
      "Status": 202
    }
  ],
  "Locale": "en",
  "InternalSignature": "",
  "CustomerID": "ef6c4c2c-b2c4-4767-987f-c298a672d4cf",
  "DeliveryService": "",
  "ShardKey": "",
  "SmID": 0,
  "DateCreated": "2026-03-03T08:47:40Z",
  "OofShard": "1"
}
```

В случае ошибки возвращается JSON с полем `error`.

---

## Веб‑интерфейс

Доступен по адресу:

```
http://localhost:8081/
```

Позволяет искать заказ по ID и отображает результат в формате JSON.

---

## Мониторинг

### Prometheus

Метрики доступны на эндпоинте:

```
/metrics
```

Основные метрики:

* `http_requests_total`
* `http_request_duration_seconds`
* `cache_hits_total`

---

### Grafana

1. Откройте `http://localhost:3000`
2. Войдите (admin / admin)
3. Добавьте Prometheus (`http://prometheus:9090`)
4. Импортируйте или создайте дашборд

<img width="989" height="806" alt="Screenshot From 2026-03-03 02-13-47" src="https://github.com/user-attachments/assets/854a23b4-bdc4-4e2e-8905-b8d589639ddc" />


---

## Схема базы данных

Таблицы:

* `orders`
* `delivery`
* `payments`
* `items`


#### ERD / схема БД
<img width="675" height="626" alt="Screenshot From 2026-03-03 11-29-20" src="https://github.com/user-attachments/assets/b975b112-5672-4619-97f1-d22803f1e686" />

---

## Нагрузочное тестирование

Для нагрузочного тестирования использовался инструмент — **Cassowary**.

Пример:

```bash
cassowary run -u http://localhost:8081/order/order-new-1 -c 10 -n 100
```

Для набора URL:

```bash
cassowary run -u http://localhost:8081 -c 5 -n 100 -f urls.txt
```
Были проведены конкурентные запросы на 100 пользователей, 100 запросов, в ходе тестирования ошибок не выявлено.
<img width="994" height="358" alt="Screenshot From 2026-03-03 11-48-12" src="https://github.com/user-attachments/assets/ea5cbbec-d703-4df0-aa02-a378c421feb5" />


## Ссылка на видеодемонстрацию
https://youtu.be/DQSUvye7Rcw
---
