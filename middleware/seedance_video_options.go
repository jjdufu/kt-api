package middleware

import (
	"encoding/json"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

// seedanceVideoTaskOptions holds the Seedance create-task fields kt-api acts
// on itself (the request body is still relayed by the Doubao plugin).
type seedanceVideoTaskOptions struct {
	ExecutionExpiresAfter json.RawMessage `json:"execution_expires_after"`
}

// SeedanceVideoTaskOptions reads gateway-side options from the Seedance create
// task body on the vendor-compatible alias route. It runs after
// PrepareTaskPluginRoute (so errors use the plugin's error renderer) and
// before Distribute and billing.
func SeedanceVideoTaskOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var options seedanceVideoTaskOptions
		if err := common.UnmarshalBodyReusable(c, &options); err != nil {
			abortTaskPluginRouteErrorDetail(c, invalidRequestStatus(err), "request body must be a JSON object")
			return
		}
		if seconds, ok := seedanceIntegerField(options.ExecutionExpiresAfter); ok {
			if normalized := service.NormalizeArkExecutionExpiresAfter(seconds); normalized > 0 {
				c.Set(service.ContextKeySeedanceExecutionExpiresAfter, normalized)
			}
		}
		c.Next()
	}
}

// seedanceIntegerField decodes an optional JSON integer; anything else is
// ignored here and left to upstream validation.
func seedanceIntegerField(raw json.RawMessage) (int64, bool) {
	if len(raw) == 0 || raw[0] == '"' {
		return 0, false
	}
	var number json.Number
	if err := common.Unmarshal(raw, &number); err != nil {
		return 0, false
	}
	value, err := number.Int64()
	if err != nil {
		return 0, false
	}
	return value, true
}
