package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/di"
	_ "github.com/tupizz/restaurant-food-golang-api-fiap/swagger"
)

// @title           FastFood Golang API
// @version         1.0
// @description     API do projeto FastFood Golang.
// @termsOfService  http://seu-site.com/terms/

// @contact.name   Suporte
// @contact.url    http://seu-site.com/support
// @contact.email  suporte@seu-site.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	container := di.BuildContainer()

	fmt.Println("Starting server")

	err := container.Invoke(func(router http.Router) {
		router.Start(":8080")
	})

	if err != nil {
		log.Fatal(err)
	}
}
