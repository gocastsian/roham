package echomiddleware

import (
	"github.com/gocastsian/roham/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AccessCheck(Roles []types.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userInfo, ok := c.Get("userInfo").(*types.UserInfo)
			if !ok || userInfo == nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "Access denied: user info not found",
				})
			}

			isAllowed := false
			for _, role := range Roles {
				if userInfo.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "Access denied: insufficient permissions",
				})
			}

			return next(c)
		}
	}
}
