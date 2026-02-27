GraphQL Comments System — это API на Go для иерархических комментариев с поддержкой вложенности любой глубины (ответы на ответы). Проект построен по принципам чистой архитектуры с разделением на слои (models, storage, gql). Реализовано два типа хранилищ: in-memory (с sync.RWMutex для потокобезопасности) и PostgreSQL (с индексами и каскадным удалением). API реализовано через GraphQL с рекурсивными типами, есть GraphiQL интерфейс для тестирования. Код покрыт unit-тестами для ключевых компонентов (storage и построение дерева комментариев). Проект запускается в Docker, используется PostgreSQL, миграции через schema.sql, управление зависимостями через Go modules.

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
