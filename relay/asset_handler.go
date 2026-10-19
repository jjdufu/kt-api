package relay

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const assetUpstreamPath = "/api/assets/"

// RelayAsset handles the 12 Seedance asset Actions as sync HTTP forwards.
func RelayAsset(c *gin.Context) {
	action := common.GetContextKeyString(c, constant.ContextKeyAssetAction)
	version := common.GetContextKeyString(c, constant.ContextKeyAssetVersion)
	if version == "" {
		version = "2024-01-01"
	}

	info, err := relaycommon.GenRelayInfo(c, types.RelayFormatOpenAI, nil, nil)
	if err != nil {
		writeAssetLocalError(c, http.StatusInternalServerError, action, version, "InternalError", err.Error())
		return
	}
	info.InitChannelMeta(c)
	if billingModel := common.GetContextKeyString(c, constant.ContextKeyAssetBillingModel); billingModel != "" {
		info.OriginModelName = billingModel
	}

	priceData, priceErr := helper.ModelPriceHelperPerCall(c, info)
	if priceErr != nil {
		writeAssetLocalError(c, http.StatusBadRequest, action, version, "InvalidParameter", priceErr.Error())
		return
	}
	info.PriceData = priceData

	session, apiErr := service.NewBillingSession(c, info, priceData.Quota)
	if apiErr != nil {
		status := apiErr.StatusCode
		if status == 0 {
			status = http.StatusForbidden
		}
		writeAssetLocalError(c, status, action, version, "InsufficientQuota", apiErr.Error())
		return
	}

	reqBody, bodyErr := prepareAssetUpstreamBody(c, info.UserId, action)
	if bodyErr != nil {
		session.Refund(c)
		writeAssetLocalError(c, http.StatusBadRequest, action, version, "InvalidParameter", bodyErr.Error())
		return
	}

	status, respBody, upstreamErr := forwardAssetRequest(c, reqBody)
	if upstreamErr != nil {
		session.Refund(c)
		writeAssetLocalError(c, http.StatusBadGateway, action, version, "UpstreamError", upstreamErr.Error())
		return
	}

	success := status >= 200 && status < 300 && !assetResponseHasError(action, respBody)
	if success {
		if applyErr := applyAssetOwnershipSideEffects(c, info, action, reqBody, respBody); applyErr != nil {
			logger.LogWarn(c, common.LogText("asset ownership side effect failed action=%s err=%s", action, applyErr.Error()))
		}
		if settleErr := session.Settle(priceData.Quota); settleErr != nil {
			logger.LogError(c, common.LogText("asset billing settle failed: %s", settleErr.Error()))
		}
		recordAssetConsumeLog(c, info, priceData.Quota, action)
		model.UpdateUserUsedQuotaAndRequestCount(info.UserId, priceData.Quota)
		model.UpdateChannelUsedQuota(info.ChannelId, priceData.Quota)
	} else {
		session.Refund(c)
	}

	c.Data(status, "application/json", respBody)
}

func readRequestBodyBytes(c *gin.Context) ([]byte, error) {
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		return nil, err
	}
	return storage.Bytes()
}

func prepareAssetUpstreamBody(c *gin.Context, userID int, action string) ([]byte, error) {
	ct := c.ContentType()
	if strings.Contains(ct, "multipart/form-data") {
		if action == "ListAssetGroups" || action == "ListAssets" {
			return nil, fmt.Errorf("%s requires application/json", action)
		}
		body, err := readRequestBodyBytes(c)
		if err != nil {
			return nil, fmt.Errorf("read multipart body failed")
		}
		return body, nil
	}

	body, err := readRequestBodyBytes(c)
	if err != nil {
		return nil, fmt.Errorf("read request body failed")
	}
	if len(body) == 0 {
		body = []byte("{}")
	}
	if action == "ListAssets" {
		return prepareListAssetsBody(userID, body)
	}
	if action != "ListAssetGroups" {
		return body, nil
	}

	var payload map[string]any
	if err := common.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	filter, _ := payload["Filter"].(map[string]any)
	if filter == nil {
		filter = map[string]any{}
		payload["Filter"] = filter
	}
	groupType, _ := filter["GroupType"].(string)
	groupType = strings.TrimSpace(groupType)
	if groupType == "" {
		return nil, fmt.Errorf("Filter.GroupType is required")
	}
	var requested []string
	if raw, ok := filter["GroupIds"].([]any); ok {
		for _, item := range raw {
			if s, ok := item.(string); ok {
				requested = append(requested, s)
			}
		}
	}
	owned, err := model.IntersectOwnedAssetGroupIDs(userID, requested, groupType, len(requested) == 0)
	if err != nil {
		return nil, err
	}
	if len(owned) == 0 {
		filter["GroupIds"] = []string{"__ktapi_no_owned_groups__"}
	} else {
		filter["GroupIds"] = owned
	}
	return common.Marshal(payload)
}

// prepareListAssetsBody makes sure every asset group a ListAssets request
// names belongs to userID. The middleware only pins the channel from one
// group, so without this check a caller could append other users' group IDs
// to Filter.GroupIds and list their assets.
func prepareListAssetsBody(userID int, body []byte) ([]byte, error) {
	var payload map[string]any
	if err := common.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	groupIDs, filterIDs, topLevel, err := collectListAssetsGroupIDs(payload)
	if err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		return nil, fmt.Errorf("Filter.GroupIds must contain at least one asset group id")
	}
	for _, id := range groupIDs {
		if _, err := model.GetAssetGroupByUser(userID, id); err != nil {
			return nil, err
		}
	}
	filter, _ := payload["Filter"].(map[string]any)
	if filter == nil {
		filter = map[string]any{}
		payload["Filter"] = filter
	}
	if len(filterIDs) == 0 {
		filterIDs = []string{topLevel}
	}
	filter["GroupIds"] = filterIDs
	return common.Marshal(payload)
}

// collectListAssetsGroupIDs returns every non-empty group id in
// Filter.GroupIds plus the top-level GroupId/group_id, de-duplicated. It also
// returns the cleaned Filter.GroupIds and the top-level id on their own.
func collectListAssetsGroupIDs(payload map[string]any) (all []string, filterIDs []string, topLevel string, err error) {
	seen := map[string]bool{}
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			all = append(all, s)
		}
	}
	for _, key := range []string{"GroupId", "group_id"} {
		raw, ok := payload[key]
		if !ok || raw == nil {
			continue
		}
		s, ok := raw.(string)
		if !ok {
			return nil, nil, "", fmt.Errorf("%s must be a string", key)
		}
		s = strings.TrimSpace(s)
		if topLevel == "" {
			topLevel = s
		}
		add(s)
	}
	if rawFilter, ok := payload["Filter"]; ok && rawFilter != nil {
		filter, ok := rawFilter.(map[string]any)
		if !ok {
			return nil, nil, "", fmt.Errorf("Filter must be an object")
		}
		if rawIDs, ok := filter["GroupIds"]; ok && rawIDs != nil {
			list, ok := rawIDs.([]any)
			if !ok {
				return nil, nil, "", fmt.Errorf("Filter.GroupIds must be an array of strings")
			}
			filterSeen := map[string]bool{}
			for _, item := range list {
				s, ok := item.(string)
				if !ok {
					return nil, nil, "", fmt.Errorf("Filter.GroupIds must be an array of strings")
				}
				s = strings.TrimSpace(s)
				if s != "" && !filterSeen[s] {
					filterSeen[s] = true
					filterIDs = append(filterIDs, s)
				}
				add(s)
			}
		}
	}
	return all, filterIDs, topLevel, nil
}

func forwardAssetRequest(c *gin.Context, body []byte) (int, []byte, error) {
	baseURL := strings.TrimRight(common.GetContextKeyString(c, constant.ContextKeyChannelBaseUrl), "/")
	if baseURL == "" {
		return 0, nil, fmt.Errorf("channel base URL is empty")
	}
	action := common.GetContextKeyString(c, constant.ContextKeyAssetAction)
	version := common.GetContextKeyString(c, constant.ContextKeyAssetVersion)
	upstreamURL := fmt.Sprintf("%s%s?Action=%s&Version=%s", baseURL, assetUpstreamPath, action, version)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	ct := c.Request.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/json"
	}
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Accept", "application/json")
	key := common.GetContextKeyString(c, constant.ContextKeyChannelKey)
	if key != "" {
		if !strings.HasPrefix(strings.ToLower(key), "bearer ") {
			key = "Bearer " + key
		}
		req.Header.Set("Authorization", key)
	}

	resp, err := service.GetHttpClient().Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, respBody, nil
}

func assetResponseHasError(action string, body []byte) bool {
	if len(body) == 0 {
		return true
	}
	switch action {
	case "CreateVisualValidateSession":
		return gjson.GetBytes(body, "BytedToken").String() == ""
	case "GetVisualValidateResult":
		return gjson.GetBytes(body, "ResponseMetadata.Error.Code").String() != ""
	default:
		return gjson.GetBytes(body, "ResponseMetadata.Error.Code").String() != ""
	}
}

func applyAssetOwnershipSideEffects(c *gin.Context, info *relaycommon.RelayInfo, action string, reqBody, respBody []byte) error {
	channelID := common.GetContextKeyInt(c, constant.ContextKeyChannelId)
	keyIndex := common.GetContextKeyInt(c, constant.ContextKeyChannelMultiKeyIndex)
	keyFP := model.HashChannelKeyFingerprint(common.GetContextKeyString(c, constant.ContextKeyChannelKey))
	userID := info.UserId

	switch action {
	case "CreateAssetGroup":
		groupID := gjson.GetBytes(respBody, "Result.Id").String()
		if groupID == "" {
			return fmt.Errorf("missing Result.Id")
		}
		name := firstNonEmpty(gjson.GetBytes(reqBody, "Name").String(), gjson.GetBytes(reqBody, "name").String())
		groupType := gjson.GetBytes(reqBody, "GroupType").String()
		if groupType == "" {
			groupType = model.AssetGroupTypeAIGC
		}
		project := gjson.GetBytes(reqBody, "ProjectName").String()
		return model.CreateAssetGroupRecord(&model.AssetGroup{
			UserId: userID, GroupId: groupID, GroupType: groupType, ProjectName: project,
			Name: name, ChannelId: channelID, KeyIndex: keyIndex, KeyFingerprint: keyFP,
		})
	case "CreateAsset":
		assetID := gjson.GetBytes(respBody, "Result.Id").String()
		if assetID == "" {
			return fmt.Errorf("missing Result.Id")
		}
		groupID := firstNonEmpty(gjson.GetBytes(reqBody, "GroupId").String(), strings.TrimSpace(c.PostForm("GroupId")))
		project := firstNonEmpty(gjson.GetBytes(reqBody, "ProjectName").String(), strings.TrimSpace(c.PostForm("ProjectName")))
		name := firstNonEmpty(gjson.GetBytes(reqBody, "Name").String(), strings.TrimSpace(c.PostForm("Name")))
		assetType := firstNonEmpty(gjson.GetBytes(reqBody, "AssetType").String(), strings.TrimSpace(c.PostForm("AssetType")))
		return model.CreateAssetRecord(&model.Asset{
			UserId: userID, AssetId: assetID, GroupId: groupID, ProjectName: project,
			Name: name, AssetType: assetType, ChannelId: channelID, KeyIndex: keyIndex, KeyFingerprint: keyFP,
		})
	case "UpdateAssetGroup":
		id := gjson.GetBytes(reqBody, "Id").String()
		var namePtr *string
		if gjson.GetBytes(reqBody, "Name").Exists() {
			n := gjson.GetBytes(reqBody, "Name").String()
			namePtr = &n
		}
		return model.UpdateAssetGroupRecord(userID, id, namePtr, gjson.GetBytes(reqBody, "Description").String())
	case "UpdateAsset":
		id := gjson.GetBytes(reqBody, "Id").String()
		return model.UpdateAssetRecord(userID, id, gjson.GetBytes(reqBody, "Name").String())
	case "DeleteAsset":
		id := gjson.GetBytes(reqBody, "Id").String()
		return model.DeleteAssetRecord(userID, id)
	case "DeleteAssetGroup":
		id := gjson.GetBytes(reqBody, "Id").String()
		return model.DeleteAssetGroupRecord(userID, id)
	case "CreateVisualValidateSession":
		token := gjson.GetBytes(respBody, "BytedToken").String()
		if token == "" {
			return fmt.Errorf("missing BytedToken")
		}
		project := gjson.GetBytes(reqBody, "ProjectName").String()
		return model.CreateVisualValidateSessionRecord(&model.VisualValidateSession{
			UserId: userID, TokenHash: model.HashBytedToken(token), ProjectName: project,
			ChannelId: channelID, KeyIndex: keyIndex, KeyFingerprint: keyFP,
		})
	case "GetVisualValidateResult":
		token := gjson.GetBytes(reqBody, "BytedToken").String()
		project := gjson.GetBytes(reqBody, "ProjectName").String()
		groupID := gjson.GetBytes(respBody, "GroupId").String()
		if groupID == "" {
			// Do not consume the session until GroupId is available.
			return fmt.Errorf("GroupId not ready")
		}
		session, err := model.ConsumeVisualValidateSession(userID, token, project)
		if err != nil {
			return err
		}
		return model.CreateAssetGroupRecord(&model.AssetGroup{
			UserId: userID, GroupId: groupID, GroupType: model.AssetGroupTypeLivenessFace,
			ProjectName: session.ProjectName, ChannelId: session.ChannelId,
			KeyIndex: session.KeyIndex, KeyFingerprint: session.KeyFingerprint,
		})
	default:
		return nil
	}
}

func recordAssetConsumeLog(c *gin.Context, info *relaycommon.RelayInfo, quota int, action string) {
	other := model.NewLogOther()
	other.SetPublic("request_path", c.Request.URL.Path)
	other.SetPublic("asset_action", action)
	other.SetPublic("model_price", info.PriceData.ModelPrice)
	other.SetPublic("group_ratio", info.PriceData.GroupRatioInfo.GroupRatio)
	model.RecordConsumeLog(c, info.UserId, model.RecordConsumeLogParams{
		ChannelId:      info.ChannelId,
		ModelName:      info.OriginModelName,
		TokenName:      c.GetString("token_name"),
		Quota:          quota,
		TokenId:        info.TokenId,
		UseTimeSeconds: int(time.Since(info.StartTime).Seconds()),
		Group:          info.UsingGroup,
		Other:          other,
		Content: []*common.Message{
			common.NewMessage("asset action: {{action}}", map[string]any{"action": action}),
		},
	})
}

func writeAssetLocalError(c *gin.Context, status int, action, version, code, message string) {
	if status == 0 {
		status = http.StatusBadRequest
	}
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
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
