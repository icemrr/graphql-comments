1. Запуск с in-memory хранилищем (без Docker)
go run ./cmd/server/main.go -storage=memory


2. Запуск с PostgreSQL
# Запустить контейнер
docker run --name postgres-comments -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d postgres:15

# Создать базу
docker exec postgres-comments createdb -U postgres comments_db

# Применить схему (PowerShell)
Get-Content schema.sql | docker exec -i postgres-comments psql -U postgres -d comments_db

# Запустить сервер
go run ./cmd/server/main.go -storage=postgres -dsn="postgres://postgres:qwerty@localhost/comments_db?sslmode=disable"

3. DataLoader еще не реализовывал.