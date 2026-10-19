package model

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupAssetModelTest(t *testing.T) {
	t.Helper()
	originalDB := DB
	t.Cleanup(func() { DB = originalDB })
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, DB.AutoMigrate(&AssetGroup{}, &Asset{}, &VisualValidateSession{}))
	// Idempotent second migrate
	require.NoError(t, DB.AutoMigrate(&AssetGroup{}, &Asset{}, &VisualValidateSession{}))
}

func TestAssetGroupOwnershipIsolation(t *testing.T) {
	setupAssetModelTest(t)

	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{
		UserId: 1, GroupId: "group-a", GroupType: AssetGroupTypeAIGC, ChannelId: 10, KeyIndex: 0,
	}))
	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{
		UserId: 2, GroupId: "group-b", GroupType: AssetGroupTypeAIGC, ChannelId: 20, KeyIndex: 1,
	}))

	got, err := GetAssetGroupByUser(1, "group-a")
	require.NoError(t, err)
	assert.Equal(t, 10, got.ChannelId)

	_, err = GetAssetGroupByUser(2, "group-a")
	assert.ErrorIs(t, err, ErrAssetNotFound)

	_, err = GetAssetGroupByUser(1, "group-b")
	assert.ErrorIs(t, err, ErrAssetNotFound)
}

func TestListAndIntersectOwnedAssetGroupIDs(t *testing.T) {
	setupAssetModelTest(t)

	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{UserId: 1, GroupId: "g1", GroupType: AssetGroupTypeAIGC, ChannelId: 1}))
	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{UserId: 1, GroupId: "g2", GroupType: AssetGroupTypeLivenessFace, ChannelId: 1}))
	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{UserId: 1, GroupId: "g3", GroupType: AssetGroupTypeAIGC, ChannelId: 1}))
	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{UserId: 9, GroupId: "other", GroupType: AssetGroupTypeAIGC, ChannelId: 9}))

	all, err := ListOwnedAssetGroupIDs(1, "")
	require.NoError(t, err)
	assert.Equal(t, []string{"g1", "g2", "g3"}, all)

	aigc, err := ListOwnedAssetGroupIDs(1, AssetGroupTypeAIGC)
	require.NoError(t, err)
	assert.Equal(t, []string{"g1", "g3"}, aigc)

	// ListAssetGroups must inject owned IDs: empty request + allowAll → all owned of type
	injected, err := IntersectOwnedAssetGroupIDs(1, nil, AssetGroupTypeAIGC, true)
	require.NoError(t, err)
	assert.Equal(t, []string{"g1", "g3"}, injected)

	// Client asks for foreign + owned → only owned remain
	filtered, err := IntersectOwnedAssetGroupIDs(1, []string{"other", "g1", "g1", "missing"}, AssetGroupTypeAIGC, false)
	require.NoError(t, err)
	assert.Equal(t, []string{"g1"}, filtered)

	empty, err := IntersectOwnedAssetGroupIDs(1, nil, AssetGroupTypeAIGC, false)
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestAssetCRUDAndCascadeDeleteGroup(t *testing.T) {
	setupAssetModelTest(t)

	require.NoError(t, CreateAssetGroupRecord(&AssetGroup{UserId: 1, GroupId: "g1", ChannelId: 5, KeyIndex: 2, KeyFingerprint: "abc"}))
	require.NoError(t, CreateAssetRecord(&Asset{
		UserId: 1, AssetId: "asset-1", GroupId: "g1", ChannelId: 5, KeyIndex: 2, KeyFingerprint: "abc", AssetType: "Image",
	}))
	require.NoError(t, CreateAssetRecord(&Asset{
		UserId: 1, AssetId: "asset-2", GroupId: "g1", ChannelId: 5, KeyIndex: 2, AssetType: "Video",
	}))

	asset, err := GetAssetByUser(1, "asset-1")
	require.NoError(t, err)
	assert.Equal(t, "g1", asset.GroupId)

	_, err = GetAssetByUser(2, "asset-1")
	assert.ErrorIs(t, err, ErrAssetNotFound)

	require.NoError(t, UpdateAssetRecord(1, "asset-1", "renamed"))
	asset, err = GetAssetByUser(1, "asset-1")
	require.NoError(t, err)
	assert.Equal(t, "renamed", asset.Name)

	require.NoError(t, DeleteAssetRecord(1, "asset-2"))
	_, err = GetAssetByUser(1, "asset-2")
	assert.ErrorIs(t, err, ErrAssetNotFound)

	require.NoError(t, DeleteAssetGroupRecord(1, "g1"))
	_, err = GetAssetGroupByUser(1, "g1")
	assert.ErrorIs(t, err, ErrAssetNotFound)
	_, err = GetAssetByUser(1, "asset-1")
	assert.ErrorIs(t, err, ErrAssetNotFound)
}

func TestEnsureAssetProject(t *testing.T) {
	assert.NoError(t, EnsureAssetProject("default", ""))
	assert.NoError(t, EnsureAssetProject("default", "default"))
	assert.ErrorIs(t, EnsureAssetProject("default", "other"), ErrAssetProjectMismatch)
}

func TestHashBytedTokenStableAndEmpty(t *testing.T) {
	assert.Empty(t, HashBytedToken(""))
	assert.Empty(t, HashBytedToken("  "))
	a := HashBytedToken("token-one")
	b := HashBytedToken("token-one")
	c := HashBytedToken("token-two")
	assert.Equal(t, a, b)
	assert.NotEqual(t, a, c)
	assert.Len(t, a, 64)
	assert.NotEqual(t, "token-one", a)
}

func TestVisualValidateSessionHashExpiryAndSingleUse(t *testing.T) {
	setupAssetModelTest(t)

	token := "byted-secret-token"
	require.NoError(t, CreateVisualValidateSessionRecord(&VisualValidateSession{
		UserId: 1, TokenHash: HashBytedToken(token), ChannelId: 7, KeyIndex: 1, ProjectName: "default",
	}))

	// Wrong user
	_, err := ConsumeVisualValidateSession(2, token, "default")
	assert.ErrorIs(t, err, ErrVisualValidateTokenInvalid)

	// Wrong project
	_, err = ConsumeVisualValidateSession(1, token, "other-project")
	assert.ErrorIs(t, err, ErrVisualValidateProjectMatch)

	// Happy path consumes once
	session, err := ConsumeVisualValidateSession(1, token, "default")
	require.NoError(t, err)
	assert.Equal(t, 7, session.ChannelId)
	assert.Equal(t, 1, session.KeyIndex)

	// Second consume fails (single-use)
	_, err = ConsumeVisualValidateSession(1, token, "default")
	assert.ErrorIs(t, err, ErrVisualValidateTokenInvalid)
}

func TestVisualValidateSessionExpiry(t *testing.T) {
	setupAssetModelTest(t)

	token := "expiring-token"
	now := common.GetTimestamp()
	require.NoError(t, CreateVisualValidateSessionRecord(&VisualValidateSession{
		UserId: 1, TokenHash: HashBytedToken(token), ChannelId: 3,
		CreatedAt: now - 100, ExpiresAt: now - 1, // already expired
	}))

	_, err := PeekVisualValidateSession(1, token)
	assert.ErrorIs(t, err, ErrVisualValidateTokenInvalid)

	_, err = ConsumeVisualValidateSession(1, token, "")
	assert.ErrorIs(t, err, ErrVisualValidateTokenInvalid)

	// Expired row should be cleaned on consume attempt
	var count int64
	require.NoError(t, DB.Model(&VisualValidateSession{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestVisualValidateSessionDefaultTTL(t *testing.T) {
	setupAssetModelTest(t)

	before := common.GetTimestamp()
	require.NoError(t, CreateVisualValidateSessionRecord(&VisualValidateSession{
		UserId: 1, TokenHash: HashBytedToken("ttl-token"), ChannelId: 1,
	}))
	after := common.GetTimestamp()

	session, err := PeekVisualValidateSession(1, "ttl-token")
	require.NoError(t, err)
	ttl := int64(VisualValidateSessionTTL / time.Second)
	assert.GreaterOrEqual(t, session.ExpiresAt, before+ttl)
	assert.LessOrEqual(t, session.ExpiresAt, after+ttl)
}
