package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupAssetIntentTest(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	originalDB := model.DB
	t.Cleanup(func() { model.DB = originalDB })
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	model.DB = db
	require.NoError(t, db.AutoMigrate(&model.Asset{}, &model.Channel{}))
}

func TestApplyAssetIntentPinsSameChannelAndKey(t *testing.T) {
	setupAssetIntentTest(t)
	base := "https://yisu.example"
	require.NoError(t, model.DB.Create(&model.Channel{
		Id: 42, Type: constant.ChannelTypeDoubaoVideo, Status: common.ChannelStatusEnabled,
		Name: "seedance", Key: "sk", BaseURL: &base,
	}).Error)
	require.NoError(t, model.CreateAssetRecord(&model.Asset{UserId: 7, AssetId: "asset-1", GroupId: "g1", ChannelId: 42, KeyIndex: 1}))
	require.NoError(t, model.CreateAssetRecord(&model.Asset{UserId: 7, AssetId: "asset-2", GroupId: "g1", ChannelId: 42, KeyIndex: 1}))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	common.SetContextKey(c, constant.ContextKeyUserId, 7)
	_ = service.GetChannelConstraints(c)

	err := applyAssetIntent(c, map[string]any{"assetIds": []any{"asset://asset-1", "asset-2"}}, pluginruntime.Meta{Key: "doubao"})
	require.Nil(t, err)
	pin, found, _ := service.GetChannelConstraints(c).ResolvedPin()
	require.True(t, found)
	assert.Equal(t, 42, pin.ChannelId)
	require.NotNil(t, pin.KeyIndex)
	assert.Equal(t, 1, *pin.KeyIndex)
	assert.Equal(t, dto.PinSourceAsset, pin.Source)
}

func TestApplyAssetIntentCrossUserDenied(t *testing.T) {
	setupAssetIntentTest(t)
	require.NoError(t, model.CreateAssetRecord(&model.Asset{UserId: 1, AssetId: "asset-x", GroupId: "g1", ChannelId: 1}))
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	common.SetContextKey(c, constant.ContextKeyUserId, 2)
	_ = service.GetChannelConstraints(c)
	err := applyAssetIntent(c, map[string]any{"assetIds": []any{"asset-x"}}, pluginruntime.Meta{})
	require.NotNil(t, err)
	assert.Equal(t, "asset_not_found", err.Code)
}

func TestPrepareAssetRouteRejectsUnknownAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/assets/", PrepareAssetRoute(), func(c *gin.Context) { c.String(200, "ok") })
	req := httptest.NewRequest(http.MethodPost, "/api/assets/?Action=Nope&Version=2024-01-01", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Action is not supported")
}
