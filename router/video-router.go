package router

import (
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetVideoRouter(router *gin.Engine) {
	videoSharedRouter := router.Group("/v1")
	videoSharedRouter.Use(middleware.RouteTag("relay"))
	videoSharedRouter.Use(middleware.TokenAuth())
	videoSharedRouter.Use(middleware.SystemPerformanceCheck())
	videoSharedRouter.POST(
		"/video/generations",
		middleware.PinTaskPluginEndpoint(),
		middleware.TaskPluginEndpointOnly(middleware.ModelRequestRateLimit()),
		middleware.PrepareTaskPluginEndpoint(),
		middleware.Distribute(),
		func(c *gin.Context) {
			controller.RelayTaskPluginEndpoint(c, controller.RelayTask)
		},
	)

	videoV1Router := router.Group("/v1")
	videoV1Router.Use(middleware.RouteTag("relay"))
	videoV1Router.Use(middleware.TokenAuth(), middleware.Distribute())
	{
		videoV1Router.GET("/video/generations/:task_id", controller.RelayTaskFetch)
		videoV1Router.POST("/videos/:video_id/remix", controller.RelayTask)
	}

	// Vendor-compatible Seedance paths (alias onto Doubao plugin routes).
	seedanceAlias := router.Group("/api/v3/contents/generations")
	seedanceAlias.Use(middleware.RouteTag("relay"))
	seedanceAlias.Use(middleware.TokenAuth())
	seedanceAlias.Use(middleware.SystemPerformanceCheck())
	{
		seedanceAlias.POST(
			"/tasks",
			middleware.PinSeedanceVideoAlias(),
			middleware.ModelRequestRateLimit(),
			middleware.PrepareTaskPluginRoute(),
			middleware.SeedanceVideoTaskOptions(),
			middleware.Distribute(),
			controller.RelayTask,
		)
		seedanceAlias.GET(
			"/tasks/:task_id",
			middleware.PinSeedanceVideoAlias(),
			middleware.PrepareTaskPluginRoute(),
		)
	}
}
