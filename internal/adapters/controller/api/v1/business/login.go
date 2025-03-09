package business

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
//	@Summary		Business login
//	@Description	Authenticate business and set secure cookie "business_auth_token" with JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.BusinessLoginRequest	true	"Business credentials"
//	@Success		200		{object}	dto.HTTPStatus
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/business/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var requestBody dto.BusinessLoginRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	business, err := h.businessService.Login(c.Request().Context(), &requestBody)
	switch {
	case errors.Is(err, errorz.ErrBusinessNotFound):
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
	cookie.Name = "business_auth_token"
	cookie.Value = business.Token
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
		Message: "Successfully logged in as business",
	})
}
