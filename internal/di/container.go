package di

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http/handler"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/repository"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Configuration
	container.Provide(config.LoadConfig)

	// Database Connection
	container.Provide(NewDatabaseConnection)

	// Router
	container.Provide(http.NewRouter)

	// Repositories
	container.Provide(repository.NewUserRepository)

	// Services
	container.Provide(service.NewUserService)

	// Handlers
	container.Provide(handler.NewUserHandler)

	return container
}
