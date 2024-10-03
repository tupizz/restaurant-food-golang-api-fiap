aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 156041436605.dkr.ecr.us-east-1.amazonaws.com
docker build -t fiap/golang-project .
docker tag fiap/golang-project:latest 156041436605.dkr.ecr.us-east-1.amazonaws.com/fiap/golang-project:latest
docker push 156041436605.dkr.ecr.us-east-1.amazonaws.com/fiap/golang-project:latest


# Introduction

Explanation:

cmd/: Entry point of the application.
internal/: Internal application code.
adapter/: Implementations of interfaces (controllers, repositories).
http/: HTTP handlers and router setup.
repository/: Data access implementations.
application/: Business logic.
service/: Application services.
dto/: Data Transfer Objects.
domain/: Core business entities and repository interfaces.
di/: Dependency Injection container setup.
