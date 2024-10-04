### creating new migration

migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" up

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable"
docker compose up -d
go run ./cmd/main.go
```

### Dependencies

**Install air for hot reload**

```bash
# binary will be /usr/local/bin/air
curl -sSfL https://goblin.run/github.com/air-verse/air | sh
```

**Install golang-migrate for migrations**

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

or 

```sh
brew install golang-migrate
```
----


```bash
docker compose up -d # start only database
air # start server with hot reload
```