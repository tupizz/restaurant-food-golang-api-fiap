package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http/handler"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(userHandler handler.UserHandler /* Add other handlers */) Router {
	engine := gin.Default()

	// Set up routes
	v1 := engine.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/", userHandler.GetAll)
			userGroup.POST("/", userHandler.Create)
		}
	}

	return Router{engine: engine}
}

func (r Router) Start(address string) {
	r.engine.Run(address)
}
