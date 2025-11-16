## Конфигурация (`.env`)

Используются переменные окружения:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=pr_service
DB_SSLMODE=disable

APP_PORT=8080
```

## 🐳 Запуск через Docker

Из каталога `deployments`:

```bash
docker-compose up --build
```

Либо через Makefile (из корня проекта):

```bash
make up
```

Остановка:

```bash
make down
# или
docker-compose down
```

## 🧹 Makefile

Основные команды:

```bash
make build      # сборка бинаря в bin/pr-service
make run        # запуск приложения локально (без Docker)
make test       # запуск go test ./...
make lint       # запуск golangci-lint
make up         # docker-compose up --build
make down       # docker-compose down
make logs       # логи контейнера app
```
---

