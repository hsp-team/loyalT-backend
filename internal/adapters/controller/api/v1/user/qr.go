package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
)

// GetQR godoc
//
//	@Summary		Get user qr code data
//	@Description	Get user qr code data
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.QRGetResponse
//	@Failure		400	{object}	dto.HTTPStatus
//	@Failure		401	{object}	dto.HTTPStatus
//	@Failure		404	{object}	dto.HTTPStatus
//	@Failure		500	{object}	dto.HTTPStatus
//	@Router			/user/qr [get]
func (h *Handler) GetQR(c echo.Context) error {
	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	response, err := h.qrService.GetUserQR(c.Request().Context(), userID)
	switch {
	case errors.Is(err, errorz.ErrUserNotFound):
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

	return c.JSON(200, response)
}
