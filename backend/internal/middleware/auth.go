package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/water-factory/api/pkg/jwt"
	"github.com/water-factory/api/pkg/errs"
)

func Auth(jwtMgr *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(401, errs.Unauthorized("missing or invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, errs.Unauthorized("invalid or expired token"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
