package api

import (
	"log"
	"net/http"

	"github.com/ecuyle/gomine/internal/servermanager"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

type ServerOptions struct {
	Name           string                 `json:"name"`
	UserID         string                 `json:"userId"`
	Runtime        string                 `json:"runtime"`
	IsEulaAccepted bool                   `json:"isEulaAccepted"`
	Config         map[string]interface{} `json:"config"`
}

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

func PostServer(context *gin.Context) {
	var options ServerOptions

	if err := context.BindJSON(&options); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	server, err := servermanager.MakeServer(options.Runtime, options.Name, options.IsEulaAccepted, options.Config)

	if err != nil {
		log.Println(err)
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	context.IndentedJSON(http.StatusCreated, server)
}

func PutServer(context *gin.Context) {
}
