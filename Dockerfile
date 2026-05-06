# 🟦 ЭТАП 1: Компиляция
FROM golang:1.25-alpine AS builder
WORKDIR /src

# Зависимости (кэшируется)
COPY go.mod go.sum ./
RUN go mod download

# Весь исходный код
COPY . .

# Сборка статического бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/server ./cmd/server

# 🟦 ЭТАП 2: Минимальный образ для запуска
FROM alpine:3.19
WORKDIR /app

# Системные сертификаты + пользователь
RUN apk --no-cache add ca-certificates && adduser -D appuser

# Бинарник
COPY --from=builder /bin/server .

# 📦 Копируем ВСЕ папки, которые ваш код читает в рантайме
# (если код читает ещё что-то: templates, config, certs и т.д. → просто добавьте COPY <папка> ./<папка>/ сюда)

COPY database/ ./database/

# Безопасность
USER appuser

# Порт берём из вашего .env (SERVER_ADDR=:8081)
EXPOSE 8081

CMD ["./server"]
