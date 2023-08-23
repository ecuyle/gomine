package main

import (
	"log"
	"os/exec"

	"github.com/ecuyle/gomine/internal/authentication"
	"github.com/ecuyle/gomine/internal/servers"
	"github.com/ecuyle/gomine/internal/token"
	"github.com/ecuyle/gomine/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cmd := exec.Command("mkdir", "-p", servers.GetJarFilepath(""), servers.GetServerFilepath(""))

	if err := cmd.Run(); err != nil {
		log.Fatalln("main.go: Could not initialize required directories")
	}

	router := gin.Default()

	serverRoutes := router.Group("/api/mcsrv")
	serverRoutes.Use(token.JwtAuthMiddleware())
	serverRoutes.GET("/", servers.GetServersByUserId)
	serverRoutes.GET("/detail", servers.GetServerDetails)
	serverRoutes.GET("/defaults", servers.GetDefaults)
	serverRoutes.POST("/", servers.PostServer)
	serverRoutes.PUT("/properties", servers.PutServerProperties)

	router.POST("/api/mcusr", user.PostUser)

	router.POST("/api/login", authentication.AuthenticateUser)

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	router.Run("localhost:8080")
}
