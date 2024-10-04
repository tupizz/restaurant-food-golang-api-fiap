# Quick Start Guide

This guide will help you set up and run the **FastFood Golang** project on your machine. Follow the steps below to get the application up and running.

## Prerequisites

- **Git**: For cloning the repository.
- **Docker** and **Docker Compose**: To containerize and run the application and database.
- **Golang** (optional): If you prefer to run the application without Docker.

## Getting Started

### 1. Clone the Repository

Open your terminal and clone the repository to your local machine:

```bash
git clone https://github.com/your-username/fastfood-golang.git
cd fastfood-golang
```

### 1.1 Install Dependencies

Ensure you have Go installed and set up on your machine.

```bash 
# binary will be /usr/local/bin/air
curl -sSfL https://goblin.run/github.com/air-verse/air | sh
````

### 2. Set Up Environment Variables

Create a `.env` file in the root directory of the project to store environment variables:

```bash
touch .env
```

Add the following content to the `.env` file:

```env
DATABASE_URL=postgres://postgres:postgres@db:5432/fiap_fast_food?sslmode=disable
```

### 3. Build and Run with Docker Compose

Use Docker Compose to build the images and start the services:

```bash
docker-compose up --build
```

This command will:

- Build the Docker image for the Go application.
- Start the PostgreSQL database container.
- Run database migrations.
- Start the Go application container.

### 4. Test the API Endpoints

Once the services are up and running, you can test the API endpoints.

#### a. Get All Users

```bash
curl http://localhost:8080/api/v1/users
```

#### b. Create a New User

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe"}' http://localhost:8080/api/v1/users
```

You should receive a JSON response with the details of the newly created user.

## Project Structure Overview

- **`cmd/your-app/main.go`**: Entry point of the application.
- **`internal/`**: Contains the application's internal code.
    - **`adapter/`**: Adapters for HTTP handlers and repository implementations.
        - **`http/`**: HTTP server setup and route definitions.
        - **`repository/`**: Database interaction implementations.
    - **`application/`**: Business logic and service definitions.
    - **`domain/`**: Core business entities and repository interfaces.
    - **`config/`**: Configuration loading and management.
    - **`di/`**: Dependency injection setup using Uber's `dig`.
- **`migrations/`**: Database migration files.
- **`Dockerfile`**: Dockerfile for building the Go application image.
- **`docker-compose.yml`**: Docker Compose configuration file.
- **`go.mod`** and **`go.sum`**: Go modules files.

## Running Without Docker (Optional)

If you prefer to run the application directly on your machine without Docker:

### 1. Install Dependencies

Ensure you have Go installed and set up on your machine.

```bash
go mod download
```

### 2. Set Up PostgreSQL

Install PostgreSQL and create a database named `yourdb`. Update the `DATABASE_URL` in the `.env` file to point to your local database:

```env
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/yourdb?sslmode=disable
```

### 3. Run Database Migrations

Install the migration tool:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run the migrations:

```bash
migrate -database ${DATABASE_URL} -path ./migrations up
```

### 4. Run the Application

```bash
go run ./cmd/your-app/main.go
```

### 5. Test the API Endpoints

Use the same `curl` commands as above to test the API.

## Troubleshooting

- **Port Conflicts**: Ensure ports `8080` (application) and `5432` (database) are not in use by other services.
- **Environment Variables**: Double-check the `DATABASE_URL` in your `.env` file.
- **Docker Issues**: If you encounter issues with Docker, try restarting Docker or running `docker-compose down` to reset the containers.

## Additional Information

- **Stopping the Services**: To stop the Docker containers, press `Ctrl+C` in the terminal where `docker-compose` is running, or run:

  ```bash
  docker-compose down
  ```

- **Rebuilding the Images**: If you make changes to the Go code or Dockerfile, rebuild the images:

  ```bash
  docker-compose up --build
  ```

- **Viewing Logs**: To view the logs of the running containers:

  ```bash
  docker-compose logs -f
  ```

- **Database Access**: You can connect to the PostgreSQL database using a tool like `psql` or any PostgreSQL client using the credentials specified in the `docker-compose.yml` file.