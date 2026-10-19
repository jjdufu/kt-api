package middleware

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

const maxAssetIDs = 32
const maxAssetIDLen = 128

type assetIntentError struct {
	Code       string
	Message    string
	StatusCode int
}

// applyAssetIntent pins the request to the channel (and key) that owns the
// referenced Seedance assets when the plugin reports assetIds in its intent.
func applyAssetIntent(c *gin.Context, intent map[string]any, _ pluginruntime.Meta) *assetIntentError {
	raw, present := intent["assetIds"]
	if !present {
		return nil
	}
	var values []any
	switch typed := raw.(type) {
	case []any:
		values = typed
	case []string:
		values = make([]any, len(typed))
		for i, id := range typed {
			values[i] = id
		}
	default:
		return &assetIntentError{Code: "invalid_asset_ids", Message: "asset ids are invalid", StatusCode: http.StatusBadRequest}
	}
	if len(values) > maxAssetIDs {
		return &assetIntentError{Code: "invalid_asset_ids", Message: "asset ids are invalid", StatusCode: http.StatusBadRequest}
	}
	ids := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		id, ok := value.(string)
		if !ok {
			return &assetIntentError{Code: "invalid_asset_ids", Message: "asset ids are invalid", StatusCode: http.StatusBadRequest}
		}
		id = strings.TrimSpace(id)
		id = strings.TrimPrefix(id, "asset://")
		if id == "" || utf8.RuneCountInString(id) > maxAssetIDLen {
			return &assetIntentError{Code: "invalid_asset_ids", Message: "asset ids are invalid", StatusCode: http.StatusBadRequest}
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil
	}

	userID := common.GetContextKeyInt(c, constant.ContextKeyUserId)
	channelID := 0
	keyIndex := 0
	for _, id := range ids {
		asset, err := model.GetAssetByUser(userID, id)
		if err != nil {
			if err == model.ErrAssetNotFound {
				return &assetIntentError{Code: "asset_not_found", Message: "asset not found or not owned by you", StatusCode: http.StatusBadRequest}
			}
			return &assetIntentError{Code: "asset_not_found", Message: "asset not found or not owned by you", StatusCode: http.StatusInternalServerError}
		}
		if channelID == 0 {
			channelID = asset.ChannelId
			keyIndex = asset.KeyIndex
		} else if asset.ChannelId != channelID || asset.KeyIndex != keyIndex {
			return &assetIntentError{Code: "asset_channel_conflict", Message: "assets must belong to the same channel", StatusCode: http.StatusBadRequest}
		}
	}

	channel, err := model.CacheGetChannel(channelID)
	if err != nil || channel == nil || channel.Status != common.ChannelStatusEnabled {
		return &assetIntentError{Code: "asset_channel_disabled", Message: "asset channel is disabled", StatusCode: http.StatusBadRequest}
	}
	idx := keyIndex
	service.GetChannelConstraints(c).AddPin(dto.ChannelPin{
		ChannelId: channel.Id,
		Source:    dto.PinSourceAsset,
		Rank:      dto.PinRankAsset,
		RetryMode: dto.PinRetrySameChannel,
		KeyIndex:  &idx,
	})
	return nil
}
