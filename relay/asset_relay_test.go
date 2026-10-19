package relay

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	kitdto "github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupAssetRelayTest(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	originalDB := model.DB
	originalLogDB := model.LOG_DB
	t.Cleanup(func() { model.DB = originalDB; model.LOG_DB = originalLogDB })
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	model.DB = db
	model.LOG_DB = db
	require.NoError(t, db.AutoMigrate(
		&model.AssetGroup{}, &model.Asset{}, &model.VisualValidateSession{},
		&model.User{}, &model.Token{}, &model.Channel{}, &model.Ability{}, &model.Log{},
		&model.UserSubscription{}, &model.SubscriptionPlan{},
	))
	ratio_setting.InitRatioSettings()
	service.InitHttpClient()
	prevRedis := common.RedisEnabled
	common.RedisEnabled = false
	t.Cleanup(func() { common.RedisEnabled = prevRedis })
}

func TestPrepareAssetUpstreamBodyInjectsOwnedGroupIDs(t *testing.T) {
	setupAssetRelayTest(t)
	require.NoError(t, model.CreateAssetGroupRecord(&model.AssetGroup{UserId: 1, GroupId: "g1", GroupType: model.AssetGroupTypeAIGC, ChannelId: 1}))
	require.NoError(t, model.CreateAssetGroupRecord(&model.AssetGroup{UserId: 1, GroupId: "g2", GroupType: model.AssetGroupTypeAIGC, ChannelId: 1}))
	require.NoError(t, model.CreateAssetGroupRecord(&model.AssetGroup{UserId: 9, GroupId: "other", GroupType: model.AssetGroupTypeAIGC, ChannelId: 9}))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := `{"Filter":{"GroupType":"AIGC","GroupIds":["other","g1"]},"PageNumber":1,"PageSize":10}`
	c.Request = httptest.NewRequest(http.MethodPost, "/api/assets/?Action=ListAssetGroups&Version=2024-01-01", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	out, err := prepareAssetUpstreamBody(c, 1, "ListAssetGroups")
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(out, &payload))
	filter := payload["Filter"].(map[string]any)
	ids := filter["GroupIds"].([]any)
	require.Len(t, ids, 1)
	assert.Equal(t, "g1", ids[0])
}

func TestPrepareAssetUpstreamBodyEmptyOwnedForcesSentinel(t *testing.T) {
	setupAssetRelayTest(t)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := `{"Filter":{"GroupType":"AIGC"},"PageNumber":1,"PageSize":10}`
	c.Request = httptest.NewRequest(http.MethodPost, "/api/assets/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	out, err := prepareAssetUpstreamBody(c, 1, "ListAssetGroups")
	require.NoError(t, err)
	assert.Contains(t, string(out), "__ktapi_no_owned_groups__")
}

func TestAssetResponseHasError(t *testing.T) {
	assert.True(t, assetResponseHasError("CreateAssetGroup", []byte(`{"ResponseMetadata":{"Error":{"Code":"X"}}}`)))
	assert.False(t, assetResponseHasError("CreateAssetGroup", []byte(`{"ResponseMetadata":{},"Result":{"Id":"g1"}}`)))
	assert.True(t, assetResponseHasError("CreateVisualValidateSession", []byte(`{}`)))
	assert.False(t, assetResponseHasError("CreateVisualValidateSession", []byte(`{"BytedToken":"abc","H5Link":"https://x"}`)))
	assert.False(t, assetResponseHasError("GetVisualValidateResult", []byte(`{}`)))
}

func TestRelayAssetCreateGroupWithMockUpstream(t *testing.T) {
	setupAssetRelayTest(t)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CreateAssetGroup", r.URL.Query().Get("Action"))
		assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ResponseMetadata":{"Action":"CreateAssetGroup","Version":"2024-01-01"},"Result":{"Id":"group-new-1"}}`)
	}))
	t.Cleanup(upstream.Close)

	user := &model.User{Id: 1, Username: "u1", Quota: 1_000_000, Status: common.UserStatusEnabled}
	user.SetSetting(kitdto.UserSetting{BillingPreference: "wallet_only"})
	require.NoError(t, model.DB.Create(user).Error)
	base := upstream.URL
	ch := &model.Channel{
		Id: 100, Type: constant.ChannelTypeDoubaoVideo, Status: common.ChannelStatusEnabled,
		Name: "yisu", Key: "sk-upstream", Models: "seedance-asset",
		BaseURL: &base, Group: "default",
	}
	require.NoError(t, model.DB.Create(ch).Error)
	require.NoError(t, model.DB.Create(&model.Ability{Group: "default", Model: "seedance-asset", ChannelId: 100, Enabled: true}).Error)
	prevCache := common.MemoryCacheEnabled
	common.MemoryCacheEnabled = true
	model.InitChannelCache()
	t.Cleanup(func() {
		common.MemoryCacheEnabled = prevCache
		model.InitChannelCache()
	})

	engine := gin.New()
	engine.POST("/api/assets/",
		func(c *gin.Context) {
			common.SetContextKey(c, constant.ContextKeyUserId, 1)
			common.SetContextKey(c, constant.ContextKeyUserQuota, 1_000_000)
			common.SetContextKey(c, constant.ContextKeyUserGroup, "default")
			common.SetContextKey(c, constant.ContextKeyUsingGroup, "default")
			common.SetContextKey(c, constant.ContextKeyTokenId, 1)
			common.SetContextKey(c, constant.ContextKeyTokenKey, "tok")
			common.SetContextKey(c, constant.ContextKeyTokenUnlimited, true)
			common.SetContextKey(c, constant.ContextKeyUserSetting, kitdto.UserSetting{BillingPreference: "wallet_only"})
			c.Set("token_name", "test")
			c.Next()
		},
		middleware.PrepareAssetRoute(),
		middleware.Distribute(),
		RelayAsset,
	)

	body := `{"Name":"demo","GroupType":"AIGC","ProjectName":"default"}`
	req := httptest.NewRequest(http.MethodPost, "/api/assets/?Action=CreateAssetGroup&Version=2024-01-01", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Contains(t, w.Body.String(), "group-new-1")

	group, err := model.GetAssetGroupByUser(1, "group-new-1")
	require.NoError(t, err)
	assert.Equal(t, 100, group.ChannelId)
}

func TestHashAndConsumeVisualValidateSessionBillingPath(t *testing.T) {
	setupAssetRelayTest(t)
	token := "session-token-xyz"
	require.NoError(t, model.CreateVisualValidateSessionRecord(&model.VisualValidateSession{
		UserId: 3, TokenHash: model.HashBytedToken(token), ChannelId: 8, KeyIndex: 0, ProjectName: "default",
	}))
	// Simulate success side-effect path used by handler
	session, err := model.ConsumeVisualValidateSession(3, token, "default")
	require.NoError(t, err)
	assert.Equal(t, 8, session.ChannelId)
	_, err = model.ConsumeVisualValidateSession(3, token, "default")
	assert.ErrorIs(t, err, model.ErrVisualValidateTokenInvalid)
}
