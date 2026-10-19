package middleware

import (
	"net/http"
	"strings"

	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"

	"github.com/gin-gonic/gin"
)

// PinSeedanceVideoAlias maps vendor-compatible /api/v3/contents/generations/tasks
// onto the Doubao plugin's declared /doubao/api/v3/... routes. Plugins cannot
// register /api/* themselves (reservedRouteNamespaces).
func PinSeedanceVideoAlias() gin.HandlerFunc {
	return func(c *gin.Context) {
		generation := pluginruntime.DefaultRegistry.Generation()
		if generation == nil {
			abortWithOpenAiMessage(c, http.StatusServiceUnavailable, "seedance video plugin is unavailable")
			return
		}
		fullPath := strings.TrimSpace(c.FullPath())
		if fullPath == "" {
			fullPath = c.Request.URL.Path
		}
		declaredPath := "/doubao" + fullPath
		binding, ok := generation.LookupDeclaredRoute(c.Request.Method, declaredPath)
		if !ok || binding.Plugin == nil {
			abortWithOpenAiMessage(c, http.StatusNotFound, "seedance video alias route is unavailable")
			return
		}
		c.Set(pluginruntime.ContextKeyPinnedPlugin, pluginruntime.PinnedPlugin{
			Generation: generation,
			Plugin:     binding.Plugin,
		})
		c.Set(pluginruntime.ContextKeyPinnedRoute, pluginruntime.PinnedRoute{
			Generation: generation,
			Plugin:     binding.Plugin,
			Route:      binding.Route,
		})
		c.Next()
	}
}
