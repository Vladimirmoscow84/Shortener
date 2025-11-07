# Shortener Service

Сервис для создания коротких ссылок с возможностью отслеживания кликов и аналитики. 

Проект написан на Go с использованием Gin, PostgreSQL и Redis.  

```


## Структура проекта
Shortener/
├── cmd/
│   └── main.go           # Точка входа приложения
├── db/
│   └── dumps/            # SQL дампы и миграции базы данных
├── internal/             # Основная логика приложения
│   ├── app/              # Запуск сервиса
│   ├── handlers/         # HTTP хэндлеры для API
│   ├── model/            # Модели данных
│   ├── service/          # Бизнес-логика
│   └── storage/          # Работа с Postgres и Redis
├── temp/                 # Временные файлы (например, фронтенд)
└── .env                  # Конфигурация среды


```

```


## Визуальное представление структуры
cmd/main.go
   │
   ▼
internal/app/Run() ──► config (.env)
        │
        ▼
storage ──► Postgres
        │
        └─► Redis (cache)
        │
        ▼
service ──► логика генерации коротких ссылок, запись кликов, аналитика
        │
        ▼
handlers ──► API эндпоинты (/shorten, /s/:short_code, /analytics/:short_url_id)
        │
        ▼
temp ──► frontend (HTML + JS)


```

``` 


## Архитектура взаимодействия сервисов
Frontend (HTML+JS)
       │
       ▼
   HTTP Requests
       │
       ▼
Handlers ──────────────► ServiceURL ──────────────► Storage
  POST /shorten         │                        ┌──► Postgres
  GET /s/:short_code    │                        └──► Redis (cache)
  GET /analytics/:id    │

Frontend: форма для создания короткой ссылки и просмотра аналитики.
Handlers: принимает HTTP запросы, валидирует их, передает в сервис.
ServiceURL: бизнес-логика:
Генерация короткой ссылки
Кэширование в Redis
Логирование кликов
Получение аналитики
Storage: работа с базой данных (Postgres) и кэшем (Redis)

```

## Установка

Склонировать репозиторий:
git clone github.com/Vladimirmoscow84/Shortener
cd Shortener

Настроить .env файл

Запустить Postgress и  Redis с нужными параметрами
docker run -d --name pg-shortener -p ________________
docker run -d --name redis-go -p _______________

 Применить  SQL миграции из db/dumpы/:
 psql -h localhost -p 5440 -U ________ -d _________ -f db/dumps/init.sql

 Запустить проект:
 go run cmd/main.go

После запуска сервис будет доступен на http://localhost:7550

## API Endpoints
Создание короткой ссылки
POST /shorten
Content-Type: application/json

{
  "original_code": "https://example.com"
}

Ответ:
{
  "id": 1,
  "short_code": "abc123XY",
  "original_code": "https://example.com",
  "created_at": "2025-11-07T14:00:00Z"
}

Переход по короткой ссылке
GET /s/:short_code
Делает редирект на оригинальный URL и сохраняет клик

Получение аналитики по короткой ссылке
GET /analytics/:short_url_id

Ответ:
{
  "2025-11-07": {
    "Mozilla/5.0": 3,
    "Chrome/120.0": 2
  }
}

## Frontend
В директории temp/ есть простой HTML + JS интерфейс для:
Создания короткой ссылки
Просмотра аналитики по ссылкам
Доступен по http://localhost:7550/

## Стек технологий
Go (Gin + sqlx + Redis)
Postgres
Redis
HTML + JS (фронтенд)
Docker (для баз)

## Тестирование

Запуск юнит-тестов:
go test ./internal/storage/...

Тестирование API через Postman или curl.


Лицензия
 © 2025 — VladimirMoscow84