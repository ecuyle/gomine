package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ecuyle/gomine/internal/servermanager"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
	"gotest.tools/assert"
)

// Test GetDefaults and assert that the expected default server properties are returned
// Path: internal/api/server.go
func TestGetDefaults(t *testing.T) {
	router := gin.Default()
	router.GET("/defaults", GetDefaults)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/defaults", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expectedProperties := servermanager.ServerProperties{}
	var p properties.Properties
	p.Decode(&expectedProperties)

	actual := servermanager.ServerProperties{}
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.DeepEqual(t, expectedProperties, actual)
}
