package routes

import (
	"medods-test-task/pkg/utils"

	_ "medods-test-task/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller interface {
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

func RegistrationRoutes(app *gin.Engine, tokenManager utils.TokenManager, c Controller) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	v1 := app.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.RefreshToken)
	}

	app.GET("/docs/*any", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
