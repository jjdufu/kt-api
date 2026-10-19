package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSelfTaskCallbackSigning(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setting := system_setting.GetTaskCallbackSetting()
	previous := *setting
	t.Cleanup(func() { *setting = previous })

	call := func() map[string]any {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/user/task_callback/signing", nil)
		c.Set("id", 42)
		GetSelfTaskCallbackSigning(c)
		require.Equal(t, http.StatusOK, recorder.Code)
		var response struct {
			Data map[string]any `json:"data"`
		}
		require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
		return response.Data
	}

	setting.SigningEnabled = false
	setting.SigningSecret = "master"
	data := call()
	assert.Equal(t, false, data["signing_enabled"])
	assert.Equal(t, "", data["secret"])

	setting.SigningEnabled = true
	data = call()
	assert.Equal(t, true, data["signing_enabled"])
	assert.Equal(t, service.TaskCallbackUserSecret("master", 42), data["secret"])

	setting.SigningSecret = ""
	data = call()
	assert.Equal(t, false, data["signing_enabled"])
}
