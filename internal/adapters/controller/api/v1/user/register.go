package user

import (
	"errors"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Register godoc
//
//	@Summary		Register new user
//	@Description	Register a new user in the system
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserRegisterRequest	true	"User registration data"
//	@Success		201		{object}	dto.UserRegisterResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/user/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var requestBody dto.UserRegisterRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	user, err := h.userService.Register(c.Request().Context(), &requestBody)
	switch {
	case errors.Is(err, errorz.ErrEmailAlreadyExists):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrConflict.Code,
		}
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	return c.JSON(201, user)
}
