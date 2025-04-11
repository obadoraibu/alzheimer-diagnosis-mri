FROM golang:1.23-alpine

RUN apk add --no-cache git postgresql-client wget tar build-base

WORKDIR /app

COPY . .

RUN go mod tidy && go mod download

RUN go build -o api ./cmd/api/main.go

# Устанавливаем postgresql-client и migrate CLI
RUN wget -q https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz \
    && tar -zxvf migrate.linux-amd64.tar.gz \
    && mv migrate.linux-amd64 /usr/bin/migrate \
    && rm -f migrate.linux-amd64.tar.gz

# Открываем порт
EXPOSE 8080

# CMD для применения миграции и запуска API
CMD ["sh", "-c", "migrate -path migrations -database postgres://$USER_DB_USER:$USER_DB_PASSWORD@$USER_DB_HOST:$USER_DB_PORT/$USER_DB_NAME?sslmode=disable up && ./api"]
