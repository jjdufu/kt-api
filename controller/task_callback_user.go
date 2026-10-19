package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

// GetSelfTaskCallbackSigning returns whether task callbacks are signed and,
// if so, the caller's signing secret for verifying X-KtApi-Signature.
func GetSelfTaskCallbackSigning(c *gin.Context) {
	setting := system_setting.GetTaskCallbackSetting()
	secret := service.CurrentTaskCallbackUserSecret(c.GetInt("id"))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"enabled":          setting.Enabled,
			"signing_enabled":  secret != "",
			"secret":           secret,
			"signature_header": service.TaskCallbackHeaderSignature,
			"timestamp_header": service.TaskCallbackHeaderTimestamp,
			"algorithm":        "v1=hex(HMAC-SHA256(secret, timestamp + \".\" + event_id + \".\" + body))",
		},
	})
}
