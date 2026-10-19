package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/service"
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
