**Lists API** — Учебный REST API для управления списками (List). Используется для практики HTTP, JSON, валидации и слоистой архитектуры. Текущее хранилище — in-memory (без БД).

•	Формат: application/json; charset=utf-8

•	Версионирование: /api/v1

•	Идентификаторы: UUID

•	Дата/время: RFC3339 (format: date-time)

•	Корреляция запросов: поддержка заголовка X-Request-Id в ответах (если передан клиентом)


Запуск:
```go mod tidy
go run ./cmd/service
# сервис будет на http://localhost:8080
```

Действия:
```
# health
curl -s http://localhost:8080/health

# создать
curl -s -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"Дом"}'

# список (с пагинацией)
curl -s "http://localhost:8080/api/v1/lists?limit=20&offset=0" -i

# получить по id
curl -s http://localhost:8080/api/v1/lists/<id>

# получить по названию
curl -s http://localhost:8080/api/v1/lists/<title>

# обновить title
curl -s -X PATCH http://localhost:8080/api/v1/lists/<id> \
  -H "Content-Type: application/json" \
  -d '{"title":"Домашние дела"}'

# удалить
curl -s -X DELETE http://localhost:8080/api/v1/lists/<id> -i
