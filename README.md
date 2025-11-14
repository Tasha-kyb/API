**Tasks API** — REST API для управления списками и задачами. Использует PostgreSQL, реализует полный CRUD цикл с валидацией и пагинацией.

**Формат:** application/json; charset=utf-8
**Версионирование:** /api/v1
**Идентификаторы:** UUID
**Дата/время:** RFC3339
**Корреляция запросов:** поддержка X-Request-Id

**Запуск:**
```bash
docker-compose up
# Или локально:
docker-compose up postgres -d
# Применить миграции
migrate -path migrations -database "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable" up
# Запустить приложение
go run ./cmd/todo-api

Примеры команд:

Работа со списками:

# 1. Создать список
curl -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" -d '{"title":"Покупки"}'

# 2. Получить все списки (с пагинацией)
curl "http://localhost:8080/api/v1/lists?limit=10&offset=0"

# 3. Получить список по ID
curl "http://localhost:8080/api/v1/lists/<list_id>"

# 4. Обновить список
curl -X PATCH http://localhost:8080/api/v1/lists/<list_id> \
  -H "Content-Type: application/json" -d '{"title":"Новое название"}'

# 5. Удалить список
curl -X DELETE "http://localhost:8080/api/v1/lists/<list_id>"

# Создать список
curl -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" -d '{"title":"Покупки"}'


Работа с задачами:

# 1. Создать задачу
curl -X POST http://localhost:8080/api/v1/lists/<list_id>/tasks \
  -H "Content-Type: application/json" -d '{"text":"Купить молоко"}'

# 2. Получить все задачи списка (с пагинацией)
curl "http://localhost:8080/api/v1/lists/<list_id>/tasks?limit=20&offset=0"

# 3. Получить задачу по ID
curl "http://localhost:8080/api/v1/tasks/<task_id>"

# 4. Обновить статус задачи (только completed)
curl -X PATCH http://localhost:8080/api/v1/tasks/<task_id> \
  -H "Content-Type: application/json" -d '{"completed":true}'

# 5. Обновить текст задачи (только text)
curl -X PATCH http://localhost:8080/api/v1/tasks/<task_id> \
  -H "Content-Type: application/json" -d '{"text":"купить тестовый тест"}'

# 6. Обновить и текст и статус (оба поля)
curl -X PATCH http://localhost:8080/api/v1/tasks/<task_id> \
  -H "Content-Type: application/json" -d '{"text":"новый текст", "completed":true}'

# 7. Удалить задачу
curl -X DELETE "http://localhost:8080/api/v1/tasks/<task_id>"


# Запустить SwaggerUI

# 1. Запустить только БД
docker-compose up -d postgres

# 2. Запустить сервер
./todo-api

# 3. Открыть Swagger
http://localhost:8080/swagger/index.html

# Тестирование кодогенерации 
go test -v ./...

# Unit тесты для сервиса
go test -v ./internal/service/

# Интеграционные тесты для репозитория
go test -v -tags=integration ./internal/storage/postgres/

# Проверка производительности индексов
docker exec -it listsapi-postgres-1 psql -U todo_user -d todo_db