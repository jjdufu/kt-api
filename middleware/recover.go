package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// PanicRecover recovers from panics in request handlers. The panic value and
// stack trace are written to the server log only; the client receives a
// generic 500 error that leaks no internal details.
func PanicRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if isClientGone(err) {
				// The client disconnected; nothing can be written back.
				common.SysLog(common.LogText("request aborted by client: %s %s: %v", c.Request.Method, c.Request.URL.Path, err))
				c.Abort()
				return
			}
			common.SysError(common.LogText("panic recovered: request_id=%s %s %s: %v\n%s",
				c.GetString(common.RequestIdKey), c.Request.Method, c.Request.URL.Path, err, debug.Stack()))
			if c.Writer.Written() {
				c.Abort()
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": "internal server error",
					"type":    "internal_error",
				},
			})
		}()
		c.Next()
	}
}

func isClientGone(recovered any) bool {
	err, ok := recovered.(error)
	if !ok {
		return false
	}
	if errors.Is(err, http.ErrAbortHandler) {
		return true
	}
	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		return false
	}
	var syscallErr *os.SyscallError
	if !errors.As(opErr, &syscallErr) {
		return false
	}
	msg := strings.ToLower(syscallErr.Error())
	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
}
