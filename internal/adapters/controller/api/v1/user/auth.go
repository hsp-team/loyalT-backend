package user

import (
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/utils"
)

func (h *Handler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_auth_token")
		if err != nil {
			return &echo.HTTPError{
				Message: "Missing authentication cookie",
				Code:    echo.ErrUnauthorized.Code,
			}
		}

		claims, err := utils.VerifyToken(cookie.Value, []byte(h.jwtConfig.UserTokenSecret()))
		if err != nil {
			return &echo.HTTPError{
				Message: err.Error(),
				Code:    echo.ErrUnauthorized.Code,
			}
		}

		c.Set("user_id", claims["sub"])

		return next(c)
	}
}
