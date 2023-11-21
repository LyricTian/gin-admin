package biz

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dal"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/cachex"
	"github.com/LyricTian/gin-admin/v10/pkg/crypto/hash"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/jwtx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login management for RBAC
type Login struct {
	Cache       cachex.Cacher
	Auth        jwtx.Auther
	UserDAL     *dal.User
	UserRoleDAL *dal.UserRole
	MenuDAL     *dal.Menu
	UserBIZ     *User
}

func (a *Login) ParseUserID(c *gin.Context) (string, error) {
	invalidToken := errors.Unauthorized(config.ErrInvalidTokenID, "Invalid access token")
	token := util.GetToken(c)
	if token == "" {
		return "", invalidToken
	}

	ctx := c.Request.Context()
	ctx = util.NewUserToken(ctx, token)

	userID, err := a.Auth.ParseSubject(ctx, token)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return "", invalidToken
		}
		return "", err
	}

	userCacheVal, ok, err := a.Cache.Get(ctx, config.CacheNSForUser, userID)
	if err != nil {
		return "", err
	} else if ok {
		userCache := util.ParseUserCache(userCacheVal)
		c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
		return userID, nil
	}

	// Check user status, if not activated, force to logout
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{SelectFields: []string{"status"}},
	})
	if err != nil {
		return "", err
	} else if user == nil || user.Status != schema.UserStatusActivated {
		return "", invalidToken
	}

	roleIDs, err := a.UserBIZ.GetRoleIDs(ctx, userID)
	if err != nil {
		return "", err
	}

	userCache := util.UserCache{
		RoleIDs: roleIDs,
	}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String())
	if err != nil {
		return "", err
	}

	c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
	return userID, nil
}

// This function generates a new captcha ID and returns it as a `schema.Captcha` struct. The length of
// the captcha is determined by the `config.C.Util.Captcha.Length` configuration value.
func (a *Login) GetCaptcha(ctx context.Context) (*schema.Captcha, error) {
	return &schema.Captcha{
		CaptchaID: captcha.NewLen(config.C.Util.Captcha.Length),
	}, nil
}

// Response captcha image
func (a *Login) ResponseCaptcha(ctx context.Context, w http.ResponseWriter, id string, reload bool) error {
	if reload && !captcha.Reload(id) {
		return errors.NotFound("", "Captcha id not found")
	}

	err := captcha.WriteImage(w, id, config.C.Util.Captcha.Width, config.C.Util.Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.NotFound("", "Captcha id not found")
		}
		return err
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *Login) genUserToken(ctx context.Context, userID string) (*schema.LoginToken, error) {
	token, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokenBuf, err := token.EncodeToJSON()
	if err != nil {
		return nil, err
	}
	logging.Context(ctx).Info("Generate user token", zap.Any("token", string(tokenBuf)))

	return &schema.LoginToken{
		AccessToken: token.GetAccessToken(),
		TokenType:   token.GetTokenType(),
		ExpiresAt:   token.GetExpiresAt(),
	}, nil
}

func (a *Login) Login(ctx context.Context, formItem *schema.LoginForm) (*schema.LoginToken, error) {
	// verify captcha
	if !captcha.VerifyString(formItem.CaptchaID, formItem.CaptchaCode) {
		return nil, errors.BadRequest(config.ErrInvalidCaptchaID, "Incorrect captcha")
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogin)

	// get user info
	user, err := a.UserDAL.GetByUsername(ctx, formItem.Username, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id", "password", "status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	} else if user.Status != schema.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	// check password
	if err := hash.CompareHashAndPassword(user.Password, formItem.Password); err != nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	}

	userID := user.ID
	ctx = logging.NewUserID(ctx, userID)

	// set user cache with role ids
	roleIDs, err := a.UserBIZ.GetRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCache := util.UserCache{RoleIDs: roleIDs}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String(),
		time.Duration(config.C.Dictionary.UserCacheExp)*time.Hour)
	if err != nil {
		logging.Context(ctx).Error("Failed to set cache", zap.Error(err))
	}
	logging.Context(ctx).Info("Login success", zap.String("username", formItem.Username))

	// generate token
	return a.genUserToken(ctx, userID)
}

func (a *Login) RefreshToken(ctx context.Context) (*schema.LoginToken, error) {
	userID := util.FromUserID(ctx)

	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest("", "Incorrect user")
	} else if user.Status != schema.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	return a.genUserToken(ctx, userID)
}

func (a *Login) Logout(ctx context.Context) error {
	userToken := util.FromUserToken(ctx)
	if userToken == "" {
		return nil
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogout)
	if err := a.Auth.DestroyToken(ctx, userToken); err != nil {
		return err
	}

	userID := util.FromUserID(ctx)
	err := a.Cache.Delete(ctx, config.CacheNSForUser, userID)
	if err != nil {
		logging.Context(ctx).Error("Failed to delete user cache", zap.Error(err))
	}
	logging.Context(ctx).Info("Logout success")

	return nil
}

// Get user info
func (a *Login) GetUserInfo(ctx context.Context) (*schema.User, error) {
	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound("", "User not found")
	}

	userRoleResult, err := a.UserRoleDAL.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	}, schema.UserRoleQueryOptions{
		JoinRole: true,
	})
	if err != nil {
		return nil, err
	}
	user.Roles = userRoleResult.Data

	return user, nil
}

// Change login password
func (a *Login) UpdatePassword(ctx context.Context, updateItem *schema.UpdateLoginPassword) error {
	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"password"},
		},
	})
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	}

	// check old password
	if err := hash.CompareHashAndPassword(user.Password, updateItem.OldPassword); err != nil {
		return errors.BadRequest("", "Incorrect old password")
	}

	// update password
	newPassword, err := hash.GeneratePassword(updateItem.NewPassword)
	if err != nil {
		return err
	}
	return a.UserDAL.UpdatePasswordByID(ctx, userID, newPassword)
}

// Query menus based on user permissions
func (a *Login) QueryMenus(ctx context.Context) (schema.Menus, error) {
	menuQueryParams := schema.MenuQueryParam{
		Status: schema.MenuStatusEnabled,
	}

	menuQueryParams.UserID = util.FromUserID(ctx)

	menuResult, err := a.MenuDAL.Query(ctx, menuQueryParams, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.MenusOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}

	// fill parent menus
	if parentIDs := menuResult.Data.SplitParentIDs(); len(parentIDs) > 0 {
		var missMenusIDs []string
		menuIDMapper := menuResult.Data.ToMap()
		for _, parentID := range parentIDs {
			if _, ok := menuIDMapper[parentID]; !ok {
				missMenusIDs = append(missMenusIDs, parentID)
			}
		}
		if len(missMenusIDs) > 0 {
			parentResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
				InIDs: missMenusIDs,
			})
			if err != nil {
				return nil, err
			}
			menuResult.Data = append(menuResult.Data, parentResult.Data...)
			sort.Sort(menuResult.Data)
		}
	}

	return menuResult.Data.ToTree(), nil
}

// Update current user info
func (a *Login) UpdateUser(ctx context.Context, updateItem *schema.UpdateCurrentUser) error {
	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID)
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	}

	user.Name = updateItem.Name
	user.Phone = updateItem.Phone
	user.Email = updateItem.Email
	user.Remark = updateItem.Remark
	return a.UserDAL.Update(ctx, user, "name", "phone", "email", "remark")
}
