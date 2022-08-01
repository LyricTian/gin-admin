package typed

import (
	"strings"
)

type Captcha struct {
	CaptchaID string `json:"captcha_id"`
}

type UserLogin struct {
	LoginName   string `json:"login_name" binding:"required"`   // login username
	Password    string `json:"password" binding:"required"`     // login password (md5)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // captcha verify id
	CaptchaCode string `json:"captcha_code" binding:"required"` // captcha verify code
}

func (a UserLogin) TrimSpace() UserLogin {
	a.LoginName = strings.TrimSpace(a.LoginName)
	a.CaptchaCode = strings.TrimSpace(a.CaptchaCode)
	return a
}

type LoginPasswordUpdate struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LoginToken struct {
	AccessToken string `json:"access_token" binding:"required"` // jwt token
	TokenType   string `json:"token_type" binding:"required"`   // Usage: (Authorization=token_type token)
	ExpiresAt   int64  `json:"expires_at" binding:"required"`   // unix timestamp
}
