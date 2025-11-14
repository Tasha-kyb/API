run:
	go run ./cmd/todo-api

build:
	go build -o bin/todo-api ./cmd/todo-api

tidy:
	go mod tidy

swagger:
    swag init -g cmd/todo-api/main.go -o ./docs

# Переменные
DB_URL=postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable

# Применить миграции
migrate-up:
    migrate -path migrations -database "$(DB_URL)" up

# Откатить одну миграцию
migrate-down:
    migrate -path migrations -database "$(DB_URL)" down 1

# Создать новую миграцию
# Использование: make migrate-create NAME=add_items_table
migrate-create:
    migrate create -ext sql -dir migrations -seq $(NAME)

# Версия миграций
migrate-version:
    migrate -path migrations -database "$(DB_URL)" version

.PHONY: migrate-up migrate-down migrate-create migrate-version