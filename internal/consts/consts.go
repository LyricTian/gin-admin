package consts

import "github.com/LyricTian/gin-admin/v10/pkg/encoding/json"

const (
	CacheNSForUser = "user:"
	CacheNSForRole = "role:"
)

const (
	// For casbin sync (if the role is deleted, the role should be deleted from the casbin)
	CacheKeyForRoleDeleted = "deleted:unix"
)

// Set user cache object
type UserCache struct {
	RoleIDs []string `json:"rids"`
}

func ParseUserCache(s string) UserCache {
	var a UserCache
	if s == "" {
		return a
	}

	_ = json.Unmarshal([]byte(s), &a)
	return a
}

func (a UserCache) String() string {
	return json.MarshalToString(a)
}
