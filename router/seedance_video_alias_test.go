package router

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSeedanceVideoAliasRoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	SetVideoRouter(engine)
	routes := engine.Routes()
	var foundPost, foundGet bool
	for _, r := range routes {
		if r.Method == http.MethodPost && r.Path == "/api/v3/contents/generations/tasks" {
			foundPost = true
		}
		if r.Method == http.MethodGet && r.Path == "/api/v3/contents/generations/tasks/:task_id" {
			foundGet = true
		}
	}
	assert.True(t, foundPost, "POST alias missing")
	assert.True(t, foundGet, "GET alias missing")
}
