package di

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/db/repository"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/http/handler"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
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

	// Repositories 2.0
	container.Provide(repository.NewClientRepository)
	container.Provide(repository.NewHealthCheckRepository)
	container.Provide(repository.NewProductRepository)
	container.Provide(repository.NewOrderRepository)
	container.Provide(repository.NewPaymentTaxSettingsRepository)
	container.Provide(repository.NewPaymentRepository)

	// UseCases
	container.Provide(usecase.NewClientUseCase)
	container.Provide(usecase.NewHealthCheckPingUseCase)
	container.Provide(usecase.NewProductUseCase)
	container.Provide(usecase.NewOrderUseCase)
	container.Provide(usecase.NewPaymentUseCase)

	// Handlers 2.0
	container.Provide(handler.NewClientHandler)
	container.Provide(handler.NewHealthcheckHandler)
	container.Provide(handler.NewProductHandler)
	container.Provide(handler.NewProductAdminHandler)
	container.Provide(handler.NewOrderHandler)
	container.Provide(handler.NewWebhookHandler)

	return container
}
