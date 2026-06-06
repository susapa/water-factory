package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/pkg/errs"
)

// permissionMap maps role → set of "module:action" strings.
var permissionMap = map[string]map[string]bool{
	"admin": {
		"master_data:read": true, "master_data:create": true, "master_data:update": true, "master_data:delete": true,
		"raw_material:read": true, "raw_material:create": true, "raw_material:update": true, "raw_material:delete": true,
		"production:read": true, "production:create": true, "production:update": true, "production:delete": true,
		"finished_goods:read": true, "finished_goods:create": true, "finished_goods:update": true, "finished_goods:delete": true,
		"sales:read": true, "sales:create": true, "sales:update": true, "sales:delete": true,
		"dashboard:read": true,
		"users:read": true, "users:create": true, "users:update": true, "users:delete": true,
	},
	"warehouse_manager": {
		"raw_material:read": true, "raw_material:create": true, "raw_material:update": true,
		"finished_goods:read": true, "finished_goods:create": true, "finished_goods:update": true,
		"production:read": true,
		"dashboard:read": true,
	},
	"production_manager": {
		"production:read": true, "production:create": true, "production:update": true,
		"raw_material:read": true,
		"dashboard:read": true,
	},
	"sales_admin": {
		"sales:read": true, "sales:create": true, "sales:update": true,
		"finished_goods:read": true,
		"dashboard:read": true,
	},
	"delivery": {
		"sales:read": true, "sales:update": true,
	},
}

func RequirePermission(module, action string) gin.HandlerFunc {
	key := module + ":" + action
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(401, errs.Unauthorized("not authenticated"))
			return
		}

		roleStr, _ := role.(string)
		if perms, ok := permissionMap[roleStr]; ok && perms[key] {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, errs.Forbidden("insufficient permissions"))
	}
}
