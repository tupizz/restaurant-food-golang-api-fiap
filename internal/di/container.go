package di

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http/handler"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/repository"
	cleanRepository "github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/db/repository"
	cleanHandler "github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/http/handler"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
	cleanUseCase "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Configuration
	container.Provide(config.LoadConfig)

	// Database Connection
	container.Provide(NewDatabaseConnection)
	container.Provide(NewSQLCDB)

	// Router
	container.Provide(http.NewRouter)

	// Repositories
	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewPaymentRepository)

	// Repositories 2.0
	container.Provide(cleanRepository.NewClientRepository)
	container.Provide(cleanRepository.NewHealthCheckRepository)
	container.Provide(cleanRepository.NewProductRepository)
	container.Provide(cleanRepository.NewOrderRepository)
	container.Provide(cleanRepository.NewPaymentTaxSettingsRepository)

	// Services
	container.Provide(service.NewUserService)
	container.Provide(service.NewPaymentService)

	// UseCases
	container.Provide(cleanUseCase.NewClientUseCase)
	container.Provide(cleanUseCase.NewHealthCheckUseCase)
	container.Provide(cleanUseCase.NewProductUseCase)
	container.Provide(cleanUseCase.NewOrderUseCase)

	// Handlers
	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewWebhookHandler)

	// Handlers 2.0
	container.Provide(cleanHandler.NewClientHandler)
	container.Provide(cleanHandler.NewHealthcheckHandler)
	container.Provide(cleanHandler.NewProductHandler)
	container.Provide(cleanHandler.NewProductAdminHandler)
	container.Provide(cleanHandler.NewOrderHandler)

	return container
}
