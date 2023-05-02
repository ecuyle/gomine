package main

import (
	"os/exec"

	"github.com/ecuyle/gomine/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	exec.Command("mkdir", "-p", "jarFiles", "worlds", "users")

	router := gin.Default()

	router.GET("/api/mcsrv", api.GetServersByUserId)
	router.GET("/api/mcsrv/defaults", api.GetDefaults)
	router.POST("/api/mcsrv", api.PostServer)
	router.PUT("/api/mcsrv", api.PutServer)

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	router.Run("localhost:8080")
}
