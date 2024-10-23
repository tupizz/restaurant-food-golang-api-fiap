# Build stage
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application's source code
COPY . .

# Build the Go application for Linux and amd64 architecture
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Run stage
FROM scratch

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port (replace with your application's port if different)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
