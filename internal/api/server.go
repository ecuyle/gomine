package api

import (
	"log"
	"net/http"

	"github.com/ecuyle/gomine/internal/db"
	"github.com/ecuyle/gomine/internal/servermanager"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

type ServersQueryParams struct {
	UserID string `json:"userId"`
}

func GetServersByUserId(context *gin.Context) {

}

func GetServerDetails(context *gin.Context) {
}

func GetDefaults(context *gin.Context) {
	var defaultProperties servermanager.ServerProperties
	var p properties.Properties

	p.Decode(&defaultProperties)
	context.IndentedJSON(http.StatusOK, &defaultProperties)
}

func respondWithInternalServerError(context *gin.Context, err error) {
	log.Println(err)
	context.String(http.StatusInternalServerError, err.Error())
}

func respondWithStatusCreated(context *gin.Context, data any) {
	context.IndentedJSON(http.StatusCreated, data)
}

func PostServer(context *gin.Context) {
	var options db.ServerOptions

	if err := context.BindJSON(&options); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	server, err := servermanager.MakeServer(options.Runtime, options.Name, options.IsEulaAccepted, options.Config)

	if err != nil {
		respondWithInternalServerError(context, err)
		return
	}

	err = db.InsertServer(server, &options)

	if err != nil {
		respondWithInternalServerError(context, err)
		return
	}

	respondWithStatusCreated(context, server)
}

func PutServer(context *gin.Context) {
}
