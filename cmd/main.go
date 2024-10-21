package main

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/di"
	_ "github.com/tupizz/restaurant-food-golang-api-fiap/swagger"
)

// @title           FastFood Golang API
// @version         1.0
// @description     API do projeto FastFood Golang.
// @termsOfService  https://tadeutupinamba.com.br

// @contact.name   Suporte
// @contact.url    https://tadeutupinamba.com.br
// @contact.email  tadeu.tupiz@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	container := di.BuildContainer()

	slog.Info("Starting server")

	err := container.Invoke(func(router http.Router) {
		slog.Info("Server started at port 8080")
		slog.Info("Swagger UI at http://localhost:8080/swagger/index.html")
		slog.Info("API Documentation at http://localhost:8080/swagger/doc.json")
		router.Start(":8080")
	})

	if err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
