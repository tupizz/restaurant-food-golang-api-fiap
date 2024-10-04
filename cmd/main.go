package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/http"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/di"
	"log"
)

func main() {
	container := di.BuildContainer()

	err := container.Invoke(func(router http.Router) {
		router.Start(":8080")
	})

	if err != nil {
		log.Fatal(err)
	}
}
