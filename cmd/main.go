package main

import (
	"medods-test-task/internal/app"
)

// @title Medods
// @version 1.0
// @description This is the test service for providing JWT
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @schemes http

func main() {
	app.Run()
}
