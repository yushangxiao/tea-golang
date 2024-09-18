package middleware

import (
	"strings"
	"tea/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/admin") {
			// 一律不使用缓存
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
			c.Header("Pragma", "no-cache")                                   // HTTP 1.0.
			c.Header("Expires", "0")                                         // Proxies.
			isAdmin := c.MustGet("userInfo").(*utils.Claims).IsAdmin
			if !isAdmin {
				c.Redirect(302, "/login")
				return
			}
			c.Next()
		}
	}
}
