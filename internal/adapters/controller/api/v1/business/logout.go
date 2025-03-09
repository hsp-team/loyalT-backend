package business

import (
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"time"
)

// Logout godoc
//
//	@Summary		Business logout
//	@Description	Delete secure cookie "business_auth_token" with JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.HTTPStatus
//	@Failure		500	{object}	dto.HTTPStatus
//	@Router			/business/auth/logout [post]
func (h *Handler) logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "business_auth_token"
	cookie.Value = "empty"
	cookie.Expires = time.Now().Add(24 * time.Hour) // Set expiration to 24 hours
	cookie.Path = "/"
	cookie.Secure = true
	if !h.devMode {
		cookie.HttpOnly = true
		cookie.SameSite = http.SameSiteStrictMode
	} else {
		cookie.SameSite = http.SameSiteNoneMode
	}
	c.SetCookie(cookie)

	return c.JSON(200, dto.HTTPStatus{
		Code:    200,
		Message: "Successfully logged out as business",
	})
}
