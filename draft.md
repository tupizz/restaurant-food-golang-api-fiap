### creating new migration

migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" up

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable"
docker compose up -d
go run ./cmd/main.go
```