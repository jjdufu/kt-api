package router

import (
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetAssetRouter(router *gin.Engine) {
	assetRouter := router.Group("/api/assets")
	assetRouter.Use(middleware.RouteTag("relay"))
	assetRouter.Use(middleware.TokenAuth())
	assetRouter.Use(middleware.SystemPerformanceCheck())
	assetRouter.Use(middleware.ModelRequestRateLimit())
	assetRouter.POST("/",
		middleware.PrepareAssetRoute(),
		middleware.Distribute(),
		controller.RelayAsset,
	)
	// Trailing-slash optional: some clients POST /api/assets?Action=...
	assetRouter.POST("",
		middleware.PrepareAssetRoute(),
		middleware.Distribute(),
		controller.RelayAsset,
	)
}
