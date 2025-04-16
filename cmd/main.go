package main

import "medods-test-task/internal/app"

// @title Medods
// @version 1.0
// @description This is the test service for providing JWT
// @host localhost
// @BasePath /v1
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @schemes http https

func main() {
	app.Run()
}
