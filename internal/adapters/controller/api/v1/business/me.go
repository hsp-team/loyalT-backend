package business

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"net/http"
)

// Me godoc
//
//	@Summary		Get business profile
//	@Description	Get business profile
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.BusinessReturn
//	@Failure		400	{object}	dto.HTTPStatus
//	@Failure		401	{object}	dto.HTTPStatus
//	@Failure		404	{object}	dto.HTTPStatus
//	@Failure		500	{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/me [get]
func (h *Handler) Me(c echo.Context) error {
	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	business, err := h.businessService.Get(c.Request().Context(), businessID)
	switch {
	case errors.Is(err, errorz.ErrBusinessNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, business)
}
