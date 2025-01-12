package http

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http/handler"
	cleanHandler "github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/http/handler"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(
	healthcheckHandler handler.HealthcheckHandler,
	userHandler handler.UserHandler,
	clientHandler cleanHandler.ClientHandler,
	productHandler handler.ProductHandler,
	adminProductHandler handler.AdminProductHandler,
	orderHandler handler.OrderHandler,
	webhookHandler handler.WebhookHandler,
) Router {
	engine := gin.Default()

	// Rota do Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Healthcheck
	engine.GET("/healthcheck", healthcheckHandler.Ping)

	// Set up routes
	v1 := engine.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/", userHandler.GetAll)
			users.POST("/", userHandler.Create)
		}

		clients := v1.Group("/clients")
		{
			clients.POST("/", clientHandler.Create)
			clients.GET("/:cpf", clientHandler.GetByCPF)
		}

		products := v1.Group("/products")
		{
			products.GET("/", productHandler.GetProducts)
		}

		orders := v1.Group("/orders")
		{
			orders.GET("/:id", orderHandler.GetById)
		}

		checkout := v1.Group("/checkout")
		{
			checkout.POST("/", orderHandler.Create)
		}

		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/notifications", webhookHandler.ProcessPayment)
		}

		admin := v1.Group("/admin")
		{
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("/", orderHandler.GetAll)
			}

			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("/", adminProductHandler.Create)
				adminProducts.PUT("/:id", adminProductHandler.Update)
				adminProducts.DELETE("/:id", adminProductHandler.Delete)
			}
		}
	}

	return Router{engine: engine}
}

func (r Router) Start(address string) {
	r.engine.Run(address)
}
