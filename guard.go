package m

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type AuthGuardPayload struct {
	RequestHeaderAuthKey string
	WhiteList            []string

	SecretKey string

	CallBack func(c *gin.Context, claims *Claims, accessToken string)
	FailBack func(c *gin.Context, code int, msg string, data any)
}

// AuthGuard http router 判断当前请求权限
//
// path 在白名单时跳过验证，可用“*”模糊匹配
func AuthGuard(p *AuthGuardPayload) gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			asterisk                    = "*"
			defaultRequestHeaderAuthKey = "Authorization"
		)

		var (
			path        = c.Request.URL.Path
			method      = c.Request.Method
			accessToken = c.GetHeader(*firstNonZeroValue[string](p.RequestHeaderAuthKey, defaultRequestHeaderAuthKey))
			claims      *Claims
			err         error
		)

		// region check symbol
		if strings.Contains(path, asterisk) {
			p.FailBack(c, 1002, fmt.Sprintf("invalid char '%s'", asterisk), nil)
			c.Abort()
			return
		}
		// endregion

		// region jump whitelist
		matchWhiteListTarget := fmt.Sprintf("%s:%s", method, path)
		for _, r := range p.WhiteList {
			if r == matchWhiteListTarget {
				return
			}
			if strings.Contains(r, asterisk) && matchPattern(strings.Split(r, ":")[1], path) {
				return
			}
		}
		// endregion

		if accessToken == "" {
			p.FailBack(c, 1001, fmt.Sprintf("未携带令牌"), nil)
			c.Abort()
			return
		}

		if claims, err = ParseToken(accessToken, p.SecretKey); err != nil {
			p.FailBack(c, 1111, fmt.Sprintf("解析AccessToken错误, %s", err.Error()), nil)
			c.Abort()
			return
		}

		if claims.ID == "" {
			p.FailBack(c, 1121, fmt.Sprintf("AccessToken无效"), nil)
			c.Abort()
			return
		}

		if time.Now().Unix() >= claims.ExpiresAt.Unix() {
			p.FailBack(c, 1121, fmt.Sprintf("AccessToken已过期"), nil)
			c.Abort()
			return
		}

		p.CallBack(c, claims, accessToken)

		c.Next()
	}
}
