package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func runSeedanceVideoTaskOptions(t *testing.T, body string) (*gin.Context, *httptest.ResponseRecorder, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v3/contents/generations/tasks", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	reached := false
	handler := SeedanceVideoTaskOptions()
	handler(c)
	if !c.IsAborted() {
		reached = true
	}
	return c, recorder, reached
}

func TestSeedanceVideoTaskOptionsExecutionExpiresAfter(t *testing.T) {
	testCases := []struct {
		name string
		body string
		want int64
	}{
		{name: "in range", body: `{"model":"m","execution_expires_after":3600}`, want: 3600},
		{name: "max", body: `{"model":"m","execution_expires_after":259200}`, want: 259200},
		{name: "below range ignored", body: `{"model":"m","execution_expires_after":60}`},
		{name: "fraction ignored", body: `{"model":"m","execution_expires_after":3600.5}`},
		{name: "string ignored", body: `{"model":"m","execution_expires_after":"3600"}`},
		{name: "absent", body: `{"model":"m"}`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, _, reached := runSeedanceVideoTaskOptions(t, testCase.body)
			assert.True(t, reached)
			value, exists := c.Get(service.ContextKeySeedanceExecutionExpiresAfter)
			if testCase.want == 0 {
				assert.False(t, exists)
				return
			}
			assert.Equal(t, testCase.want, value)
		})
	}
}

func TestSeedanceVideoTaskOptionsRejectsNonObjectBody(t *testing.T) {
	_, recorder, reached := runSeedanceVideoTaskOptions(t, `[1,2]`)
	assert.False(t, reached)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "request body must be a JSON object")
}

func TestSeedanceVideoTaskOptionsCallbackURL(t *testing.T) {
	setting := system_setting.GetTaskCallbackSetting()
	previous := *setting
	t.Cleanup(func() { *setting = previous })
	setting.Enabled = true
	setting.AllowPrivateIp = false
	setting.HTTPSOnly = false

	t.Run("valid public IP target", func(t *testing.T) {
		c, _, reached := runSeedanceVideoTaskOptions(t, `{"model":"m","callback_url":" https://93.184.216.34/cb?k=v "}`)
		assert.True(t, reached)
		target := service.SeedanceCallbackTargetFromContext(c)
		if assert.NotNil(t, target) {
			assert.Equal(t, "https://93.184.216.34/cb?k=v", target.URL)
			assert.Equal(t, "93.184.216.34", target.Host)
		}
	})
	for name, body := range map[string]string{
		"absent": `{"model":"m"}`,
		"null":   `{"model":"m","callback_url":null}`,
		"empty":  `{"model":"m","callback_url":"   "}`,
	} {
		t.Run(name, func(t *testing.T) {
			c, _, reached := runSeedanceVideoTaskOptions(t, body)
			assert.True(t, reached)
			assert.Nil(t, service.SeedanceCallbackTargetFromContext(c))
		})
	}
	for name, testCase := range map[string]struct{ body, message string }{
		"number":       {`{"model":"m","callback_url":123}`, "callback_url must be a string"},
		"object":       {`{"model":"m","callback_url":{"u":1}}`, "callback_url must be a string"},
		"private ip":   {`{"model":"m","callback_url":"http://192.168.50.51:8080/cb"}`, "callback_url is not allowed"},
		"metadata ip":  {`{"model":"m","callback_url":"http://169.254.169.254/"}`, "callback_url is not allowed"},
		"bad scheme":   {`{"model":"m","callback_url":"ftp://93.184.216.34/"}`, "callback_url must use http or https"},
		"numeric host": {`{"model":"m","callback_url":"http://2130706433/"}`, "canonical"},
	} {
		t.Run(name, func(t *testing.T) {
			c, recorder, reached := runSeedanceVideoTaskOptions(t, testCase.body)
			assert.False(t, reached)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.Contains(t, recorder.Body.String(), testCase.message)
			assert.Nil(t, service.SeedanceCallbackTargetFromContext(c))
		})
	}

	t.Run("https only", func(t *testing.T) {
		setting.HTTPSOnly = true
		defer func() { setting.HTTPSOnly = false }()
		_, recorder, reached := runSeedanceVideoTaskOptions(t, `{"model":"m","callback_url":"http://93.184.216.34/cb"}`)
		assert.False(t, reached)
		assert.Contains(t, recorder.Body.String(), "callback_url must use https")
	})

	t.Run("allow private for testing", func(t *testing.T) {
		setting.AllowPrivateIp = true
		defer func() { setting.AllowPrivateIp = false }()
		c, _, reached := runSeedanceVideoTaskOptions(t, `{"model":"m","callback_url":"http://192.168.50.51:8080/cb"}`)
		assert.True(t, reached)
		assert.NotNil(t, service.SeedanceCallbackTargetFromContext(c))
	})

	t.Run("disabled ignores callback_url", func(t *testing.T) {
		setting.Enabled = false
		defer func() { setting.Enabled = true }()
		c, _, reached := runSeedanceVideoTaskOptions(t, `{"model":"m","callback_url":"http://127.0.0.1/"}`)
		assert.True(t, reached)
		assert.Nil(t, service.SeedanceCallbackTargetFromContext(c))
	})
}
