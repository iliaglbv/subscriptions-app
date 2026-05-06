# subscriptions-app
Для запуска потребуется
1 Go (Golang) — версия 1.21 или выше
2 PostgreSQL — база данных
---
Необходимо подключится к серверу postgreSQL и создать бд subscriptions_db
В .env нужно указать пароль владельца бд
пример postgres
DATABASE_DSN=postgres://postgres:Пароль_postgres@localhost:5432/subscriptions_db?sslmode=disable
