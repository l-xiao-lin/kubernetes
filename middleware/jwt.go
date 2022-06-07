package middleware

import (
	"github.com/gin-gonic/gin"
	"lubanKubernets/utils"
	"strings"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/k8s/login" || c.Request.URL.Path == "/api/k8s/register" {
			c.Next()
		} else {
			authHeader := c.Request.Header.Get("Authorization")
			if authHeader == "" {
				c.JSON(200, gin.H{
					"msg": "请求头中auth为空",
				})
				c.Abort()
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(200, gin.H{
					"msg": "请求头的auth格式有误",
				})
				c.Abort()
				return
			}

			claims, err := utils.ParseToken(parts[1])
			if err != nil {
				c.JSON(200, gin.H{
					"msg": "无效的token",
				})
				c.Abort()
				return
			}
			c.Set("username", claims.Username)
			c.Next()
		}

	}

}
