package controller

import (
	"io"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

const taskUpstreamCallbackMaxBody = 1 << 20

// HandleTaskUpstreamCallback receives upstream task callbacks for the
// optional fast path. The body is never trusted: a valid per-task token
// only triggers an immediate re-query through the polling path, and the
// handler answers at once so upstream's 5-second window is met.
func HandleTaskUpstreamCallback(c *gin.Context) {
	setting := system_setting.GetTaskCallbackSetting()
	if !setting.Enabled || setting.UpstreamCallbackBase() == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "not found"}})
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, taskUpstreamCallbackMaxBody))
	if err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": gin.H{"code": "body_too_large", "message": "request body too large"}})
		return
	}
	var hint struct {
		Status string `json:"status"`
	}
	_ = common.Unmarshal(body, &hint)
	if len(hint.Status) > 20 {
		hint.Status = hint.Status[:20]
	}
	if !service.HandleTaskUpstreamCallback(c.Param("token"), hint.Status) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
