package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

const (
	AssetAPIVersion = "2024-01-01"

	AssetModelManagement      = "seedance-asset"
	AssetModelVisualValidate  = "seedance-visual-validate"
	AssetModelVisualValidateR = "seedance-visual-validate-result"
)

var assetActionWhitelist = map[string]string{
	"CreateAssetGroup":            AssetModelManagement,
	"CreateAsset":                 AssetModelManagement,
	"ListAssetGroups":             AssetModelManagement,
	"GetAssetGroup":               AssetModelManagement,
	"ListAssets":                  AssetModelManagement,
	"GetAsset":                    AssetModelManagement,
	"UpdateAssetGroup":            AssetModelManagement,
	"UpdateAsset":                 AssetModelManagement,
	"DeleteAsset":                 AssetModelManagement,
	"DeleteAssetGroup":            AssetModelManagement,
	"CreateVisualValidateSession": AssetModelVisualValidate,
	"GetVisualValidateResult":     AssetModelVisualValidateR,
}

// PrepareAssetRoute validates Action/Version, sets the virtual billing model,
// and pins the channel when the request references an owned group/asset/session.
func PrepareAssetRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := strings.TrimSpace(c.Query("Action"))
		version := strings.TrimSpace(c.Query("Version"))
		if version == "" {
			version = AssetAPIVersion
		}
		billingModel, ok := assetActionWhitelist[action]
		if !ok || action == "" {
			abortWithAssetError(c, http.StatusBadRequest, action, version, "InvalidParameter", "Action is not supported")
			return
		}
		if version != AssetAPIVersion {
			abortWithAssetError(c, http.StatusBadRequest, action, version, "InvalidParameter", "Version must be "+AssetAPIVersion)
			return
		}

		common.SetContextKey(c, constant.ContextKeyAssetAction, action)
		common.SetContextKey(c, constant.ContextKeyAssetVersion, version)
		common.SetContextKey(c, constant.ContextKeyAssetBillingModel, billingModel)
		common.SetContextKey(c, constant.ContextKeyOriginalModel, billingModel)
		c.Set("original_model", billingModel)

		if err := applyAssetOwnershipPin(c, action); err != nil {
			code, message := assetPinError(err)
			abortWithAssetError(c, http.StatusBadRequest, action, version, code, message)
			return
		}
		c.Next()
	}
}

func applyAssetOwnershipPin(c *gin.Context, action string) error {
	userID := common.GetContextKeyInt(c, constant.ContextKeyUserId)
	switch action {
	case "CreateAssetGroup", "CreateVisualValidateSession", "ListAssetGroups":
		return nil
	case "CreateAsset":
		groupID, err := peekAssetStringField(c, "GroupId", "group_id")
		if err != nil {
			return err
		}
		if groupID == "" {
			return fmt.Errorf("GroupId is required")
		}
		return pinFromAssetGroup(c, userID, groupID)
	case "ListAssets":
		groupIDs, err := peekListAssetsGroupIDs(c)
		if err != nil {
			return err
		}
		if len(groupIDs) == 0 {
			return fmt.Errorf("Filter.GroupIds must contain at least one asset group id")
		}
		// Every requested group must belong to the caller; the channel is
		// pinned from the first one. The relay handler re-checks the body
		// before forwarding it upstream.
		for _, groupID := range groupIDs[1:] {
			if _, err := model.GetAssetGroupByUser(userID, groupID); err != nil {
				return err
			}
		}
		return pinFromAssetGroup(c, userID, groupIDs[0])
	case "GetAssetGroup", "UpdateAssetGroup", "DeleteAssetGroup":
		groupID, err := peekAssetStringField(c, "Id", "id")
		if err != nil {
			return err
		}
		if groupID == "" {
			return fmt.Errorf("Id is required")
		}
		return pinFromAssetGroup(c, userID, groupID)
	case "GetAsset", "UpdateAsset", "DeleteAsset":
		assetID, err := peekAssetStringField(c, "Id", "id")
		if err != nil {
			return err
		}
		if assetID == "" {
			return fmt.Errorf("Id is required")
		}
		return pinFromAsset(c, userID, assetID)
	case "GetVisualValidateResult":
		token, err := peekAssetStringField(c, "BytedToken", "byted_token")
		if err != nil {
			return err
		}
		if token == "" {
			return model.ErrVisualValidateTokenInvalid
		}
		session, err := model.PeekVisualValidateSession(userID, token)
		if err != nil {
			return err
		}
		return addAssetChannelPin(c, session.ChannelId, session.KeyIndex)
	default:
		return nil
	}
}

func pinFromAssetGroup(c *gin.Context, userID int, groupID string) error {
	group, err := model.GetAssetGroupByUser(userID, groupID)
	if err != nil {
		return err
	}
	return addAssetChannelPin(c, group.ChannelId, group.KeyIndex)
}

func pinFromAsset(c *gin.Context, userID int, assetID string) error {
	asset, err := model.GetAssetByUser(userID, assetID)
	if err != nil {
		return err
	}
	return addAssetChannelPin(c, asset.ChannelId, asset.KeyIndex)
}

func addAssetChannelPin(c *gin.Context, channelID int, keyIndex int) error {
	channel, err := model.CacheGetChannel(channelID)
	if err != nil || channel == nil || channel.Status != common.ChannelStatusEnabled {
		return fmt.Errorf("asset channel is disabled")
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

func assetPinError(err error) (string, string) {
	if err == nil {
		return "InvalidParameter", "invalid request"
	}
	switch {
	case err == model.ErrAssetNotFound:
		return "InvalidParameter", "record not found"
	case err == model.ErrVisualValidateTokenInvalid:
		return "InvalidParameter", err.Error()
	case err == model.ErrVisualValidateProjectMatch:
		return "InvalidParameter", err.Error()
	case err == model.ErrAssetProjectMismatch:
		return "InvalidParameter", err.Error()
	default:
		msg := err.Error()
		if strings.Contains(msg, "disabled") {
			return "InvalidParameter", msg
		}
		return "InvalidParameter", msg
	}
}

func abortWithAssetError(c *gin.Context, status int, action, version, code, message string) {
	c.JSON(status, gin.H{
		"ResponseMetadata": gin.H{
			"RequestId": c.GetString(common.RequestIdKey),
			"Action":    action,
			"Version":   version,
			"Service":   "ark",
			"Region":    "cn-beijing",
			"Error": gin.H{
				"Code":    code,
				"Message": message,
			},
		},
	})
	c.Abort()
}

// peekAssetStringField reads a top-level JSON string or multipart form field
// without consuming the body for later forwarding (uses UnmarshalBodyReusable / PostForm).
func peekAssetStringField(c *gin.Context, keys ...string) (string, error) {
	ct := c.ContentType()
	if strings.Contains(ct, "multipart/form-data") || strings.Contains(ct, "application/x-www-form-urlencoded") {
		for _, key := range keys {
			if v := strings.TrimSpace(c.PostForm(key)); v != "" {
				return v, nil
			}
		}
		return "", nil
	}
	var body map[string]any
	if err := common.UnmarshalBodyReusable(c, &body); err != nil {
		return "", fmt.Errorf("invalid JSON body")
	}
	for _, key := range keys {
		if raw, ok := body[key]; ok {
			if s, ok := raw.(string); ok {
				return strings.TrimSpace(s), nil
			}
		}
	}
	return "", nil
}

// peekListAssetsGroupIDs returns the top-level GroupId/group_id followed by
// every non-empty Filter.GroupIds entry, de-duplicated, without consuming the
// body. ListAssets is JSON only.
func peekListAssetsGroupIDs(c *gin.Context) ([]string, error) {
	var body map[string]any
	if err := common.UnmarshalBodyReusable(c, &body); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	seen := map[string]bool{}
	var ids []string
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			ids = append(ids, s)
		}
	}
	for _, key := range []string{"GroupId", "group_id"} {
		if s, ok := body[key].(string); ok {
			add(s)
		}
	}
	filter, _ := body["Filter"].(map[string]any)
	if filter != nil {
		raw, _ := filter["GroupIds"].([]any)
		for _, item := range raw {
			if s, ok := item.(string); ok {
				add(s)
			}
		}
	}
	return ids, nil
}
