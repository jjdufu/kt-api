package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// panicLogFormat keeps the literal "panic detected" in every log language so
// the KTAPIPanic log alert keeps matching.
const panicLogFormat = "panic detected: request_id=%s %s %s: %v\n%s"

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
			// Gin-style marker through the rotation-aware writer, so the
			// current log file always records the recovery.
			_, _ = fmt.Fprintf(common.ErrorLogWriter, "[Recovery] %s panic recovered: %v\n", time.Now().Format("2006/01/02 - 15:04:05"), err)
			common.SysError(common.LogText(panicLogFormat,
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
