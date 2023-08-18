package main

import (
	"log"
	"os/exec"

	"github.com/ecuyle/gomine/internal/api"
	"github.com/ecuyle/gomine/internal/servermanager"
	"github.com/gin-gonic/gin"
)

func main() {
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

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	router.Run("localhost:8080")
}
