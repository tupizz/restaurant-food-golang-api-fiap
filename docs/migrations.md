Sure! Below is a guide in Markdown format that you can use as documentation for setting up and using `golang-migrate` for managing migrations in your Go project.

# Setting Up `golang-migrate` for Database Migrations

This guide will help you set up and use `golang-migrate` in your Go project to manage your database schema migrations.

## Prerequisites

- Go installed (v1.16 or later)
- Database server (e.g., PostgreSQL, MySQL) running and accessible

## Step 1: Install `golang-migrate` CLI

To use `golang-migrate`, you first need to install the CLI tool. There are a couple of ways to install it:

### Option 1: Install via Go

Use the following command to install `golang-migrate`:

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

This will install the `migrate` command in your Go binaries directory (`$HOME/go/bin` by default).

### Option 2: Install via Homebrew (macOS/Linux)

If you are using macOS or Linux, you can install `golang-migrate` using Homebrew:

```sh
brew install golang-migrate
```

## Step 2: Add `golang-migrate` to Your PATH (Optional)

If you installed `golang-migrate` via Go, ensure the Go binaries directory is in your system's PATH.

1. Edit your `.zshrc` or `.bashrc` file:
   ```sh
   nano ~/.zshrc
   ```

2. Add the following line to include Go binaries in your PATH:
   ```sh
   export PATH=$PATH:$HOME/go/bin
   ```

3. Save the file and reload your profile:
   ```sh
   source ~/.zshrc
   ```

4. Verify the installation:
   ```sh
   migrate -version
   ```

## Step 3: Create Migration Files

Use the `migrate` command to create new migration files.

```sh
migrate create -ext sql -dir migrations -seq create_users_table
```

This command will generate two files:

- `000001_create_users_table.up.sql`: Defines the schema changes to be applied.
- `000001_create_users_table.down.sql`: Defines the schema changes to be reverted.

You can edit these files to define your database changes. For example:

**`migrations/000001_create_users_table.up.sql`**:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**`migrations/000001_create_users_table.down.sql`**:
```sql
DROP TABLE IF EXISTS users;
```

## Step 4: Run Migrations

To apply your migrations, use the following command:

```sh
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" up
```

- `-path ./migrations` specifies the directory where your migration files are located.
- `-database` specifies the database connection string.

### Rolling Back Migrations

To roll back the most recent migration:

```sh
migrate -path ./migrations -database  "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" down 1
```

## Step 5: Using `golang-migrate` in Code (Optional)

If you need to run migrations programmatically from within your Go application:

1. Add `golang-migrate` to your project:

   ```sh
   go get -u github.com/golang-migrate/migrate/v4
   ```

2. Write code to run migrations:

   ```go
   package main

   import (
       "log"

       "github.com/golang-migrate/migrate/v4"
       _ "github.com/golang-migrate/migrate/v4/database/postgres"
       _ "github.com/golang-migrate/migrate/v4/source/file"
   )

   func main() {
       m, err := migrate.New(
           "file://migrations",
           "postgres://user:password@localhost:5432/dbname?sslmode=disable",
       )
       if err != nil {
           log.Fatal(err)
       }

       if err := m.Up(); err != nil && err != migrate.ErrNoChange {
           log.Fatal(err)
       }

       log.Println("Migrations ran successfully")
   }
   ```

## Summary

- **Installation**: Install `golang-migrate` using Go or Homebrew.
- **Creating Migrations**: Use `migrate create` to generate up and down SQL scripts.
- **Running Migrations**: Use the `migrate` command with appropriate options to apply or roll back migrations.
- **Using in Code**: You can integrate `golang-migrate` with Go code for programmatically managing migrations.

## Common Commands

- **Create a new migration**:
  ```sh
  migrate create -ext sql -dir migrations -seq migration_name
  ```

- **Apply all migrations**:
  ```sh
  migrate -path ./migrations -database "DB_CONNECTION_STRING" up
  ```

- **Roll back last migration**:
  ```sh
  migrate -path ./migrations -database "DB_CONNECTION_STRING" down 1
  ```

## Troubleshooting

- **Command Not Found**: Ensure `migrate` is installed and added to your PATH.
- **Database Connection Issues**: Check your connection string and ensure the database server is running.

For more detailed documentation, visit the official `golang-migrate` GitHub repository: [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)
```

You can save this guide in your repository as `MIGRATION_GUIDE.md` or something similar, so team members or contributors can easily follow the steps to set up and use `golang-migrate`.