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
	container.Provide(NewSQLCDB)

	// Router
	container.Provide(http.NewRouter)

	// Repositories
	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewClientRepository)
	container.Provide(repository.NewProductRepository)
	container.Provide(repository.NewOrderRepository)
	container.Provide(repository.NewPaymentTaxSettingsRepository)
	container.Provide(repository.NewPaymentRepository)

	// Services
	container.Provide(service.NewUserService)
	container.Provide(service.NewClientService)
	container.Provide(service.NewProductService)
	container.Provide(service.NewProductServiceAdmin)
	container.Provide(service.NewOrderService)
	container.Provide(service.NewPaymentService)

	// Handlers
	container.Provide(handler.NewHealthcheckHandler)
	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewClientHandler)
	container.Provide(handler.NewProductHandler)
	container.Provide(handler.NewAdminProductHandler)
	container.Provide(handler.NewOrderHandler)
	container.Provide(handler.NewFakeCheckoutHandler)

	return container
}
