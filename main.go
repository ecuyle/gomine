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

	router.GET("/api/mcsrv", api.GetServersByUserId)
	router.GET("/api/mcsrv/detail", api.GetServerDetails)
	router.GET("/api/mcsrv/defaults", api.GetDefaults)
	router.POST("/api/mcsrv", api.PostServer)
	router.PUT("/api/mcsrv/properties", api.PutServerProperties)

	router.POST("/api/mcusr", api.PostUser)

	router.POST("/api/login", api.AuthenticateUser)

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	router.Run("localhost:8080")
}
