FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api ./cmd/api/main.go

RUN apk add --no-cache postgresql-client
RUN wget -q https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz && \
    tar -zxvf migrate.linux-amd64.tar.gz && \
    mv migrate.linux-amd64 /usr/bin/migrate && \
    rm -f migrate.linux-amd64.tar.gz


EXPOSE 8080

CMD ["sh", "-c", "migrate -path migrations -database postgres://$USER_DB_USER:$USER_DB_PASSWORD@$USER_DB_HOST:$USER_DB_PORT/$USER_DB_NAME?sslmode=disable up && ./api"]