package main

import (
	"log"
	"os/exec"

	"github.com/ecuyle/gomine/internal/api"
	"github.com/ecuyle/gomine/internal/servermanager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cmd := exec.Command("mkdir", "-p", servermanager.GetJarFilepath(""), servermanager.GetServerFilepath(""))

	if err := cmd.Run(); err != nil {
		log.Fatalln("main.go: Could not initialize required directories")
	}

	router := gin.Default()

	serverRoutes := router.Group("/api/mcsrv")
	serverRoutes.Use(api.JwtAuthMiddleware())
	serverRoutes.GET("/", api.GetServersByUserId)
	serverRoutes.GET("/detail", api.GetServerDetails)
	serverRoutes.GET("/defaults", api.GetDefaults)
	serverRoutes.POST("/", api.PostServer)
	serverRoutes.PUT("/properties", api.PutServerProperties)

	router.POST("/api/mcusr", api.PostUser)

	router.POST("/api/login", api.AuthenticateUser)

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	router.Run("localhost:8080")
}
