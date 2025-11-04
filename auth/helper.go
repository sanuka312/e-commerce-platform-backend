package auth

import (
	"e-commerce-platform-backend/data"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) *data.IntrospectResponse {
	if claims, exists := c.Get("claims"); exists {
		return claims.(*data.IntrospectResponse)
	}
	return nil
}

func GetResourceRoles(claims data.ResourceRoles) []string {
	var roles []string
	for _, v := range claims {
		roles = append(roles, v.Roles...)
	}
	return roles
}

func CheckRole(roles []string, role string) bool {
	if roles == nil {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func GetRolesFromContext(c *gin.Context) []string {
	if roles, exists := c.Get("roles"); exists {
		return roles.([]string)
	}
	return nil
}

func GetUserNameFromContext(c *gin.Context) string {
	if username, exists := c.Get("user_name"); exists {
		return username.(string)
	}
	return ""
}

func IsAuthenticated(c *gin.Context) bool {
	claims := GetClaims(c)
	return claims != nil && claims.Active
}
