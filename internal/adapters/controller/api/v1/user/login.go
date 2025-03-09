package user

import (
	"errors"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"time"
)

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and set secure cookie "user_auth_token" with JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserLoginRequest	true	"User credentials"
//	@Success		200		{object}	dto.HTTPStatus
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/user/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var requestBody dto.UserLoginRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	user, err := h.userService.Login(c.Request().Context(), &requestBody)
	switch {
	case errors.Is(err, errorz.ErrUserNotFound):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrNotFound.Code,
		}
	case errors.Is(err, errorz.ErrPasswordDoesNotMatch):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrForbidden.Code,
		}
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	// Set secure cookie with JWT token
	cookie := new(http.Cookie)
	cookie.Name = "user_auth_token"
	cookie.Value = user.Token
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

	// Return success without exposing the token
	return c.JSON(200, dto.HTTPStatus{
		Code:    200,
		Message: "Successfully logged in as user",
	})
}
