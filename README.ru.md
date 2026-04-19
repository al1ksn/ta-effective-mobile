# Сервис управления подписками

> [English](README.md) | [Русский](README.ru.md)

REST API сервис для управления пользовательскими подписками, написанный на Go.

## Возможности

- Полный CRUD для подписок
- Подсчёт суммарной стоимости с фильтрацией по периоду, пользователю и сервису
- Автоматические миграции базы данных при запуске
- Swagger UI для интерактивного просмотра API
- Docker-окружение с PostgreSQL

## Технологии

- **Go 1.26** — язык
- **Chi** — HTTP роутер
- **PostgreSQL 16** — база данных
- **pgx v5** — драйвер PostgreSQL
- **golang-migrate** — миграции БД
- **Swagger / swaggo** — документация API

## Быстрый старт

### Требования

- [Docker](https://docs.docker.com/get-docker/) и Docker Compose

### Запуск через Docker

```bash
cp .env.example .env
docker-compose up
```

API будет доступно по адресу `http://localhost:8080`.  
Swagger UI: `http://localhost:8080/swagger/`

### Локальный запуск

```bash
# Установить зависимости
go mod download

# Настроить окружение
cp .env.example .env
# Отредактировать .env, указав данные PostgreSQL

# Запустить сервер
go run ./cmd/api
```

## Конфигурация

| Переменная    | По умолчанию   | Описание                  |
|---------------|----------------|---------------------------|
| `DB_HOST`     | `localhost`    | Хост PostgreSQL           |
| `DB_PORT`     | `5432`         | Порт PostgreSQL           |
| `DB_USER`     | `postgres`     | Пользователь БД           |
| `DB_PASSWORD` | `postgres`     | Пароль БД                 |
| `DB_NAME`     | `subscriptions`| Имя базы данных           |
| `SERVER_PORT` | `8080`         | Порт HTTP сервера         |

## Справочник API

Базовый путь: `/api/v1`

| Метод    | Эндпоинт                  | Описание                              |
|----------|---------------------------|---------------------------------------|
| `POST`   | `/subscriptions`          | Создать подписку                      |
| `GET`    | `/subscriptions`          | Получить список всех подписок         |
| `GET`    | `/subscriptions/{id}`     | Получить подписку по ID               |
| `PUT`    | `/subscriptions/{id}`     | Обновить подписку                     |
| `DELETE` | `/subscriptions/{id}`     | Удалить подписку                      |
| `GET`    | `/subscriptions/total`    | Подсчитать суммарную стоимость        |

### Создание / обновление подписки

```json
{
  "service_name": "Netflix",
  "price": 990,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2026",
  "end_date": "12-2026"
}
```

> Даты передаются в формате `MM-YYYY`. Поле `end_date` необязательно.

### Параметры запроса суммарной стоимости

| Параметр       | Тип    | Описание                              |
|----------------|--------|---------------------------------------|
| `from`         | string | Начало периода в формате `MM-YYYY`    |
| `to`           | string | Конец периода в формате `MM-YYYY`     |
| `user_id`      | UUID   | Фильтр по пользователю                |
| `service_name` | string | Фильтр по названию сервиса            |

### Пример ответа

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Netflix",
  "price": 990,
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "start_date": "01-2026",
  "end_date": "12-2026",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

## Схема базы данных

```sql
CREATE TABLE subscriptions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price        INTEGER      NOT NULL CHECK (price > 0),
    user_id      UUID         NOT NULL,
    start_date   DATE         NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMPTZ  DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  DEFAULT NOW()
);
```

## Структура проекта

```
.
├── cmd/api/          # Точка входа
├── internal/
│   ├── config/       # Конфигурация окружения
│   ├── handler/      # HTTP обработчики
│   ├── model/        # Модели данных
│   └── repository/   # Слой работы с БД
├── migrations/       # SQL миграции
├── docs/             # Сгенерированная Swagger документация
├── Dockerfile
└── docker-compose.yml
```

## Лицензия

MIT © 2026 Alexander Kalugin
