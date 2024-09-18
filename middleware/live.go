package middleware

import (
	"strings"
	"tea/datacon/user"
	"tea/utils"

	"github.com/gin-gonic/gin"
)

func LiveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedPaths := []string{"/api/login", "/api/signup", "/login", "/signup", "/back.jpg"}
		// 放行不需要授权的路径
		requestPath := c.Request.URL.Path
		for _, path := range allowedPaths {
			if path == requestPath || strings.HasPrefix(requestPath, path) {
				c.Next()
				return
			}
		}
		if strings.HasPrefix(requestPath, "/layer") {
			c.Next()
			return
		}
		userinfo, live := c.Get("userInfo")
		if !live {
			c.Redirect(302, "/login")
			return
		}
		userName := userinfo.(*utils.Claims).Username
		_, err := user.MyUserManager.GetUserInfo(userName)
		if err != nil {
			c.Redirect(302, "/login")
			// c.Abort()
			return
		}
		c.Next()
	}
}
