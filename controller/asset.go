package controller

import (
	"github.com/QuantumNous/new-api/relay"

	"github.com/gin-gonic/gin"
)

// RelayAsset forwards Seedance asset Actions through the native asset relay path.
func RelayAsset(c *gin.Context) {
	relay.RelayAsset(c)
}
