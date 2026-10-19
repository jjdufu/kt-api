package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPanicRecoverReturnsGenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(PanicRecover())
	router.GET("/boom", func(c *gin.Context) {
		panic("secret internal detail")
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/boom", nil))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	var body struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, "internal server error", body.Error.Message)
	require.Equal(t, "internal_error", body.Error.Type)
	require.NotContains(t, recorder.Body.String(), "secret internal detail")
	require.NotContains(t, recorder.Body.String(), "github.com")
}

func TestPanicRecoverKeepsPartialResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(PanicRecover())
	router.GET("/stream", func(c *gin.Context) {
		c.String(http.StatusOK, "partial")
		panic("after write")
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/stream", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "partial", recorder.Body.String())
}

func TestPanicRecoverAbortHandlerIsSilent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(PanicRecover())
	router.GET("/abort", func(c *gin.Context) {
		panic(http.ErrAbortHandler)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/abort", nil))

	require.Empty(t, recorder.Body.String())
}
