package model

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

const (
	AssetGroupTypeAIGC         = "AIGC"
	AssetGroupTypeLivenessFace = "LivenessFace"

	VisualValidateSessionTTL = 30 * time.Minute
)

var (
	ErrAssetNotFound              = errors.New("record not found")
	ErrAssetProjectMismatch       = errors.New("project name does not match")
	ErrVisualValidateTokenInvalid = errors.New("BytedToken is missing, expired, used, or does not belong to current user")
	ErrVisualValidateProjectMatch = errors.New("project name does not match the one used when creating the session")
)

// AssetGroup maps a local user to an upstream Seedance/Ark asset group.
// ChannelId and KeyIndex pin later asset and video calls to the same supplier credential.
type AssetGroup struct {
	Id             int64  `json:"id" gorm:"primaryKey"`
	CreatedAt      int64  `json:"created_at" gorm:"index"`
	UpdatedAt      int64  `json:"updated_at"`
	UserId         int    `json:"user_id" gorm:"uniqueIndex:uk_asset_groups_user_group;index"`
	GroupId        string `json:"group_id" gorm:"type:varchar(128);uniqueIndex:uk_asset_groups_user_group;index"`
	GroupType      string `json:"group_type" gorm:"type:varchar(32);index"`
	ProjectName    string `json:"project_name" gorm:"type:varchar(128);default:default"`
	Name           string `json:"name" gorm:"type:varchar(128)"`
	ChannelId      int    `json:"channel_id" gorm:"index"`
	KeyIndex       int    `json:"key_index" gorm:"default:0"`
	KeyFingerprint string `json:"-" gorm:"type:varchar(64)"`
}

func (AssetGroup) TableName() string { return "asset_groups" }

// Asset maps a local user to an upstream Seedance/Ark asset.
type Asset struct {
	Id             int64  `json:"id" gorm:"primaryKey"`
	CreatedAt      int64  `json:"created_at" gorm:"index"`
	UpdatedAt      int64  `json:"updated_at"`
	UserId         int    `json:"user_id" gorm:"uniqueIndex:uk_assets_user_asset;index"`
	AssetId        string `json:"asset_id" gorm:"type:varchar(128);uniqueIndex:uk_assets_user_asset;index"`
	GroupId        string `json:"group_id" gorm:"type:varchar(128);index"`
	ProjectName    string `json:"project_name" gorm:"type:varchar(128);default:default"`
	Name           string `json:"name" gorm:"type:varchar(128)"`
	AssetType      string `json:"asset_type" gorm:"type:varchar(32)"`
	ChannelId      int    `json:"channel_id" gorm:"index"`
	KeyIndex       int    `json:"key_index" gorm:"default:0"`
	KeyFingerprint string `json:"-" gorm:"type:varchar(64)"`
}

func (Asset) TableName() string { return "assets" }

// VisualValidateSession stores only the SHA-256 of BytedToken for routing
// GetVisualValidateResult back to the channel used at session create.
type VisualValidateSession struct {
	Id             int64  `json:"id" gorm:"primaryKey"`
	CreatedAt      int64  `json:"created_at" gorm:"index"`
	UserId         int    `json:"user_id" gorm:"index"`
	TokenHash      string `json:"-" gorm:"type:varchar(64);uniqueIndex"`
	ProjectName    string `json:"project_name" gorm:"type:varchar(128);default:default"`
	ChannelId      int    `json:"channel_id" gorm:"index"`
	KeyIndex       int    `json:"key_index" gorm:"default:0"`
	KeyFingerprint string `json:"-" gorm:"type:varchar(64)"`
	ExpiresAt      int64  `json:"expires_at" gorm:"index"`
}

func (VisualValidateSession) TableName() string { return "visual_validate_sessions" }

// HashBytedToken returns the hex SHA-256 of a BytedToken. Empty input yields "".
func HashBytedToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", digest)
}

// HashChannelKeyFingerprint hashes a supplier API key for multi-key pinning.
func HashChannelKeyFingerprint(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", digest)
}

func normalizeProjectName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "default"
	}
	return name
}

func touchAssetTimestamps(createdAt *int64, updatedAt *int64) {
	now := common.GetTimestamp()
	if createdAt != nil && *createdAt == 0 {
		*createdAt = now
	}
	if updatedAt != nil {
		*updatedAt = now
	}
}

// CreateAssetGroupRecord registers an upstream group for the owning user.
func CreateAssetGroupRecord(group *AssetGroup) error {
	if group == nil || group.UserId == 0 || strings.TrimSpace(group.GroupId) == "" {
		return fmt.Errorf("invalid asset group")
	}
	group.GroupId = strings.TrimSpace(group.GroupId)
	group.ProjectName = normalizeProjectName(group.ProjectName)
	if group.GroupType == "" {
		group.GroupType = AssetGroupTypeAIGC
	}
	touchAssetTimestamps(&group.CreatedAt, &group.UpdatedAt)
	return DB.Create(group).Error
}

// GetAssetGroupByUser returns a group owned by userId, or ErrAssetNotFound.
func GetAssetGroupByUser(userId int, groupId string) (*AssetGroup, error) {
	groupId = strings.TrimSpace(groupId)
	if userId == 0 || groupId == "" {
		return nil, ErrAssetNotFound
	}
	var group AssetGroup
	err := DB.Where("user_id = ? AND group_id = ?", userId, groupId).First(&group).Error
	exist, err := RecordExist(err)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrAssetNotFound
	}
	return &group, nil
}

// ListOwnedAssetGroupIDs returns all upstream group IDs for userId.
// When groupType is non-empty, results are filtered by GroupType.
func ListOwnedAssetGroupIDs(userId int, groupType string) ([]string, error) {
	if userId == 0 {
		return nil, nil
	}
	q := DB.Model(&AssetGroup{}).Where("user_id = ?", userId)
	if strings.TrimSpace(groupType) != "" {
		q = q.Where("group_type = ?", strings.TrimSpace(groupType))
	}
	var ids []string
	if err := q.Order("id asc").Pluck("group_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// IntersectOwnedAssetGroupIDs keeps only IDs the user owns.
// requested nil/empty means "all owned" when allowAll is true; otherwise empty.
func IntersectOwnedAssetGroupIDs(userId int, requested []string, groupType string, allowAll bool) ([]string, error) {
	owned, err := ListOwnedAssetGroupIDs(userId, groupType)
	if err != nil {
		return nil, err
	}
	if len(requested) == 0 {
		if allowAll {
			return owned, nil
		}
		return nil, nil
	}
	ownedSet := make(map[string]struct{}, len(owned))
	for _, id := range owned {
		ownedSet[id] = struct{}{}
	}
	out := make([]string, 0, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, id := range requested {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := ownedSet[id]; !ok {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

// UpdateAssetGroupRecord updates mutable local fields after a successful upstream update.
func UpdateAssetGroupRecord(userId int, groupId string, name *string, descriptionIgnored string) error {
	group, err := GetAssetGroupByUser(userId, groupId)
	if err != nil {
		return err
	}
	updates := map[string]any{"updated_at": common.GetTimestamp()}
	if name != nil {
		updates["name"] = strings.TrimSpace(*name)
	}
	_ = descriptionIgnored // description is upstream-only; local table keeps name for admin convenience
	return DB.Model(group).Updates(updates).Error
}

// DeleteAssetGroupRecord removes the group and all local assets under it.
// Call only after upstream delete succeeds (or when intentionally cleaning orphans).
func DeleteAssetGroupRecord(userId int, groupId string) error {
	groupId = strings.TrimSpace(groupId)
	if userId == 0 || groupId == "" {
		return ErrAssetNotFound
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND group_id = ?", userId, groupId).Delete(&AssetGroup{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrAssetNotFound
		}
		if err := tx.Where("user_id = ? AND group_id = ?", userId, groupId).Delete(&Asset{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateAssetRecord registers an upstream asset for the owning user.
func CreateAssetRecord(asset *Asset) error {
	if asset == nil || asset.UserId == 0 || strings.TrimSpace(asset.AssetId) == "" || strings.TrimSpace(asset.GroupId) == "" {
		return fmt.Errorf("invalid asset")
	}
	asset.AssetId = strings.TrimSpace(asset.AssetId)
	asset.GroupId = strings.TrimSpace(asset.GroupId)
	asset.ProjectName = normalizeProjectName(asset.ProjectName)
	touchAssetTimestamps(&asset.CreatedAt, &asset.UpdatedAt)
	return DB.Create(asset).Error
}

// GetAssetByUser returns an asset owned by userId, or ErrAssetNotFound.
func GetAssetByUser(userId int, assetId string) (*Asset, error) {
	assetId = strings.TrimSpace(assetId)
	if userId == 0 || assetId == "" {
		return nil, ErrAssetNotFound
	}
	var asset Asset
	err := DB.Where("user_id = ? AND asset_id = ?", userId, assetId).First(&asset).Error
	exist, err := RecordExist(err)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrAssetNotFound
	}
	return &asset, nil
}

// UpdateAssetRecord updates the local name after a successful upstream update.
func UpdateAssetRecord(userId int, assetId string, name string) error {
	asset, err := GetAssetByUser(userId, assetId)
	if err != nil {
		return err
	}
	return DB.Model(asset).Updates(map[string]any{
		"name":       strings.TrimSpace(name),
		"updated_at": common.GetTimestamp(),
	}).Error
}

// DeleteAssetRecord removes a local asset row after upstream delete succeeds.
func DeleteAssetRecord(userId int, assetId string) error {
	assetId = strings.TrimSpace(assetId)
	if userId == 0 || assetId == "" {
		return ErrAssetNotFound
	}
	res := DB.Where("user_id = ? AND asset_id = ?", userId, assetId).Delete(&Asset{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAssetNotFound
	}
	return nil
}

// EnsureAssetProject matches stored project name when the caller supplies one.
func EnsureAssetProject(stored string, requested string) error {
	if strings.TrimSpace(requested) == "" {
		return nil
	}
	if normalizeProjectName(stored) != normalizeProjectName(requested) {
		return ErrAssetProjectMismatch
	}
	return nil
}

// CreateVisualValidateSessionRecord stores a hashed BytedToken with TTL.
func CreateVisualValidateSessionRecord(session *VisualValidateSession) error {
	if session == nil || session.UserId == 0 || session.TokenHash == "" || session.ChannelId == 0 {
		return fmt.Errorf("invalid visual validate session")
	}
	session.ProjectName = normalizeProjectName(session.ProjectName)
	now := common.GetTimestamp()
	if session.CreatedAt == 0 {
		session.CreatedAt = now
	}
	if session.ExpiresAt == 0 {
		session.ExpiresAt = now + int64(VisualValidateSessionTTL.Seconds())
	}
	return DB.Create(session).Error
}

// ConsumeVisualValidateSession looks up a still-valid session by plaintext token
// for the given user, deletes it (single-use), and returns the row.
func ConsumeVisualValidateSession(userId int, bytedToken string, projectName string) (*VisualValidateSession, error) {
	hash := HashBytedToken(bytedToken)
	if userId == 0 || hash == "" {
		return nil, ErrVisualValidateTokenInvalid
	}
	var session VisualValidateSession
	err := DB.Where("user_id = ? AND token_hash = ?", userId, hash).First(&session).Error
	exist, err := RecordExist(err)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrVisualValidateTokenInvalid
	}
	now := common.GetTimestamp()
	if session.ExpiresAt > 0 && now > session.ExpiresAt {
		_ = DB.Delete(&session).Error
		return nil, ErrVisualValidateTokenInvalid
	}
	if err := EnsureAssetProject(session.ProjectName, projectName); err != nil {
		return nil, ErrVisualValidateProjectMatch
	}
	res := DB.Where("id = ? AND user_id = ?", session.Id, userId).Delete(&VisualValidateSession{})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrVisualValidateTokenInvalid
	}
	return &session, nil
}

// PeekVisualValidateSession returns a non-consumed session without deleting it (tests / diagnostics).
func PeekVisualValidateSession(userId int, bytedToken string) (*VisualValidateSession, error) {
	hash := HashBytedToken(bytedToken)
	if userId == 0 || hash == "" {
		return nil, ErrVisualValidateTokenInvalid
	}
	var session VisualValidateSession
	err := DB.Where("user_id = ? AND token_hash = ?", userId, hash).First(&session).Error
	exist, err := RecordExist(err)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrVisualValidateTokenInvalid
	}
	if session.ExpiresAt > 0 && common.GetTimestamp() > session.ExpiresAt {
		return nil, ErrVisualValidateTokenInvalid
	}
	return &session, nil
}
