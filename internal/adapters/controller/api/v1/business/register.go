package business

import (
	"errors"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Register godoc
//
//	@Summary		Register new business
//	@Description	Register a new business in the system
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.BusinessRegisterRequest	true	"Business registration data"
//	@Success		201		{object}	dto.BusinessRegisterResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/business/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var requestBody dto.BusinessRegisterRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	business, err := h.businessService.Register(c.Request().Context(), &requestBody)
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

	return c.JSON(201, business)
}
