package middleware

import (
	"strings"
	"tea/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey []byte

func CheckJWT(tokenString string) (*utils.Claims, error) {
	config := utils.MyConfig
	jwtKey = []byte(config.JwtKey)
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
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
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			c.Redirect(302, "/login")
			// c.Abort()
			return
		}
		tokenString := cookie.Value
		userInfo, err := CheckJWT(tokenString)
		if err != nil || tokenString == "" {
			c.Redirect(302, "/login")
			return
		}

		c.Set("userInfo", userInfo)
		c.Next()
	}
}
