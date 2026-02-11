package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"graphql-comments/internal/gql"
	"graphql-comments/internal/storage"
)

func main() {
	// Определяем флаги командной строки
	storageType := flag.String("storage", "memory", "Тип хранилища: memory или postgres")
	dsn := flag.String("dsn", "postgres://postgres:qwerty@localhost/comments_db?sslmode=disable", "DSN для PostgreSQL")
	port := flag.String("port", "8081", "Порт для HTTP сервера")
	flag.Parse()

	fmt.Println("Запуск GraphQL сервера")

	var store storage.Storage
	var err error

	// Выбор реализации хранилища
	switch *storageType {
	case "memory":
		// In-memory хранилище (данные в оперативной памяти)
		store = storage.NewMemoryStorage()
		fmt.Println("Используется in-memory хранилище")

	case "postgres":
		// PostgreSQL хранилище (данные в базе данных)
		store, err = storage.NewPostgresStorage(*dsn)
		if err != nil {
			log.Fatal("Ошибка подключения к PostgreSQL:", err)
		}
		// Приведение типа чтобы вызвать Close() только для PostgresStorage
		if pgStorage, ok := store.(*storage.PostgresStorage); ok {
			defer pgStorage.Close()
		}
		fmt.Println("Используется PostgreSQL хранилище")

	default:
		log.Fatal("Ошибка. Используйте: memory или postgres")
	}

	// Создаем GraphQL схему с переданным хранилищем
	schema, err := gql.BuildSchema(store)
	if err != nil {
		log.Fatal("Ошибка создания GraphQL схемы:", err)
	}

	// Создаем HTTP handler для GraphQL с включенным GraphiQL
	http.Handle("/graphql", gql.NewHandler(schema))

	// Запускаем HTTP сервер
	addr := ":" + *port
	fmt.Printf("Сервер запущен: http://localhost%s/graphql\n", addr)
	fmt.Println("\n  Параметры запуска:")
	fmt.Printf("   Хранилище: %s\n", *storageType)
	if *storageType == "postgres" {
		fmt.Printf("   PostgreSQL DSN: %s\n", *dsn)
	}
	fmt.Printf("   Порт: %s\n", *port)

	// Запускаем сервер (блокирующий вызов)
	log.Fatal(http.ListenAndServe(addr, nil))
}
