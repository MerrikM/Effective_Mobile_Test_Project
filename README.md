# Subscription Service API

## Описание проекта

API для управления подписками пользователей, включая создание, получение, обновление и удаление подписок, а также подсчёт общей стоимости подписок с фильтрацией по пользователю, сервису и датам. Сервис написан на Go и использует PostgreSQL в качестве базы данных. Документация API генерируется с помощью `godoc` и доступна через Swagger по адресу: `http://localhost:8080/swagger/index.html`.

## Требования

- **Go**: 1.24
- **PostgreSQL**: 16
- **Docker** и **Docker Compose** (для запуска в контейнерах)
- **Swagger**: для просмотра документации API
- Библиотеки:
    - `github.com/jmoiron/sqlx` - для работы с БД 
    - `github.com/go-chi/chi/v5` — для маршрутизации
    - `github.com/swaggo/http-swagger` — для Swagger-документации
    - Другие зависимости указаны в `go.mod`

## Установка и запуск

### Локальная установка

1. **Клонируйте репозиторий**:
   ```bash
   git clone <URL_репозитория>
   cd subscription-service
   ```

2. **Установите зависимости**:
   ```bash
   go mod tidy
   ```

3. **Настройте конфигурацию**:
    - Создайте или отредактируйте файл `config.yaml` в корне проекта. Пример:
      ```yaml
      databaseConfig:
        dsn: "postgresql://postgres:postgres@localhost:5432/mydb?sslmode=disable"
      serverAddr: ":8080"
      ```

4. **Настройте базу данных**:
    - Убедитесь, что PostgreSQL запущен.
    - Создайте базу данных `mydb` и выполните SQL-скрипт `init/init_database.sql` для создания необходимых таблиц.

5. **Запустите приложение**:
   ```bash
   go run main.go
   ```

6. **Доступ к API**:
    - API будет доступен по адресу `http://localhost:8080`.
    - Документация Swagger: `http://localhost:8080/swagger/index.html`.

### Запуск с Docker

1. **Соберите и запустите контейнеры**:
   ```bash
   docker-compose up --build
   ```

2. **Доступ к API**:
    - API: `http://localhost:8080`
    - Swagger: `http://localhost:8080/swagger/index.html`
    - База данных: `localhost:5432` (пользователь: `postgres`, пароль: `postgres`, база: `mydb`)

3. **Остановка**:
   ```bash
   docker-compose down
   ```

## API Эндпоинты

### Создание подписки
- **Эндпоинт**: `POST /subscriptions/create`
- **Описание**: Создаёт новую подписку.
- **Тело запроса**:
  ```json
  {
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-01-2025",
    "end_date": "10-12-2027"
  }
  ```
- **Успешный ответ (201)**:
  ```json
  {
    "message": "подписка успешно создана",
    "subscription": {
      "service_name": "Yandex Plus",
      "price": 400,
      "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
      "start_date": "07-01-2025",
      "end_date": "10-12-2027"
    }
  }
  ```
- **Ошибки**:
    - `400`: Неверный формат запроса
    - `500`: Ошибка создания подписки

### Получение подписок пользователя
- **Эндпоинт**: `GET /subscriptions/user/{uuid}`
- **Описание**: Возвращает список подписок по UUID пользователя.
- **Параметры**:
    - `uuid` (path): UUID пользователя
- **Успешный ответ (200)**:
  ```json
  [
    {
      "service_name": "Yandex Plus",
      "price": 400,
      "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
      "start_date": "07-01-2025",
      "end_date": "10-12-2027"
    }
  ]
  ```
- **Ошибки**:
    - `404`: Не удалось получить подписки

### Получение подписки по ID
- **Эндпоинт**: `GET /subscriptions/get/{id}`
- **Описание**: Возвращает подписку по её ID.
- **Параметры**:
    - `id` (path): ID подписки
- **Успешный ответ (200)**:
  ```json
  {
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-01-2025",
    "end_date": "10-12-2027"
  }
  ```
- **Ошибки**:
    - `400`: Неверный ID
    - `404`: Не удалось получить подписку

### Получение общей стоимости подписок
- **Эндпоинт**: `GET /subscriptions/total-cost`
- **Описание**: Возвращает сумму стоимости подписок с фильтрацией.
- **Параметры**:
    - `user_id` (query, обязательно): UUID пользователя
    - `service_name` (query, опционально): Название сервиса
    - `start_date` (query, опционально): Дата начала (DD-MM-YYYY, по умолчанию 01-01-2000)
    - `end_date` (query, опционально): Дата окончания (DD-MM-YYYY)
- **Успешный ответ (200)**:
  ```json
  {
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "общая_стоимость": 1200
  }
  ```
- **Ошибки**:
    - `400`: Ошибка параметров запроса
    - `500`: Ошибка сервера

### Обновление подписки
- **Эндпоинт**: `PUT /subscriptions/update/{id}`
- **Описание**: Обновляет подписку по ID.
- **Параметры**:
    - `id` (path): ID подписки
- **Тело запроса**:
  ```json
  {
    "service_name": "Yandex Plus",
    "price": 500,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-01-2025",
    "end_date": "10-12-2028"
  }
  ```
- **Успешный ответ (200)**:
  ```json
  {
    "message": "подписка успешно обновлена"
  }
  ```
- **Ошибки**:
    - `400`: Неверный ID или формат запроса
    - `500`: Не удалось обновить подписку

### Удаление подписки
- **Эндпоинт**: `DELETE /subscriptions/delete/{id}`
- **Описание**: Удаляет подписку по ID.
- **Параметры**:
    - `id` (path): ID подписки
- **Успешный ответ (204)**: Нет тела ответа
- **Ошибки**:
    - `400`: Неверный ID
    - `500`: Не удалось удалить подписку

## Swagger Документация

Интерактивная документация API доступна через Swagger UI по адресу: `http://localhost:8080/swagger/index.html`. Здесь вы можете:
- Просматривать описания всех эндпоинтов.
- Тестировать запросы непосредственно в браузере.
- Ознакомиться с форматами запросов и ответов, включая примеры и схемы данных.

Для генерации Swagger-документации используется пакет `github.com/swaggo/http-swagger`. Убедитесь, что аннотации `@title`, `@version`, `@description`, и т.д. в коде актуальны перед запуском.

## Структура данных

### CreateUpdateSubscriptionRequest
```go
type CreateUpdateSubscriptionRequest struct {
    ServiceName string              `json:"service_name" example:"Yandex Plus" description:"Название сервиса подписки"`
    Price       int                 `json:"price" example:"400" description:"Цена подписки"`
    UserID      string              `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" description:"UUID пользователя"`
    StartDate   model.DayMonthYear  `json:"start_date" example:"07-01-2025" description:"Дата начала подписки"`
    EndDate     *model.DayMonthYear `json:"end_date,omitempty" example:"10-12-2027" description:"Дата окончания подписки"`
}
```

### SubscriptionCreateResponse
```go
type SubscriptionCreateResponse struct {
    Message      string                    `json:"message" example:"подписка успешно создана"`
    Subscription model.SubscriptionDetails `json:"subscription"`
}
```

### TotalCostResponse
```go
type TotalCostResponse struct {
    UserID    string `json:"user_id"`
    TotalCost int    `json:"общая_стоимость"`
}
```

### SubscriptionUpdateResponse
```go
type SubscriptionUpdateResponse struct {
    Message string `json:"message" example:"подписка успешно обновлена"`
}
```

## Тестирование

1. **Локальное тестирование**:
    - Используйте инструменты вроде `curl` или Postman для отправки запросов.
    - Пример запроса для создания подписки:
      ```bash
      curl -X POST http://localhost:8080/subscriptions/create \
      -H "Content-Type: application/json" \
      -d '{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-01-2025","end_date":"10-12-2027"}'
      ```
    - Тестируйте эндпоинты через Swagger UI: `http://localhost:8080/swagger/index.html`.

2. **Тесты в коде**:
    - Добавьте модульные тесты в директорию `tests/`, используя `testing` пакет Go.

## Дополнительная информация

- **Логирование**: Ошибки логируются с помощью `log.Println`.
- **Конфигурация**: Параметры сервера и базы данных задаются в `config.yaml`.
- **Swagger**: Документация API автоматически генерируется и доступна по `/swagger/index.html`.