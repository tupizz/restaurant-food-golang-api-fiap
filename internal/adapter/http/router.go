package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http/handler"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(
	userHandler handler.UserHandler,
	clientHandler handler.ClientHandler,
) Router {
	engine := gin.Default()

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
	}

	return Router{engine: engine}
}

func (r Router) Start(address string) {
	r.engine.Run(address)
}
