package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithInternalServerError(context *gin.Context, err error) {
	log.Println(err)
	context.String(http.StatusInternalServerError, err.Error())
}

func RespondWithStatusCreated(context *gin.Context, data any) {
	context.IndentedJSON(http.StatusCreated, data)
}

func RespondWithNotFound(context *gin.Context, err error) {
	log.Println(err)
	context.String(http.StatusNotFound, err.Error())
}

func RespondWithStatusOk(context *gin.Context, data any) {
	context.IndentedJSON(http.StatusOK, data)
}
