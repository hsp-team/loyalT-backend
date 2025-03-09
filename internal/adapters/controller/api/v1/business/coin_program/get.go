package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"net/http"
)

// Get godoc
//
//	@Summary		Get coin program
//	@Description	Get business coin program. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.CoinProgramReturn
//	@Failure		400	{object}	dto.HTTPStatus
//	@Failure		401	{object}	dto.HTTPStatus
//	@Failure		403	{object}	dto.HTTPStatus
//	@Failure		404	{object}	dto.HTTPStatus
//	@Failure		500	{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program [get]
func (h *Handler) Get(c echo.Context) error {
	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	coinProgram, err := h.coinProgramService.Get(c.Request().Context(), businessID)
	switch {
	case errors.Is(err, errorz.ErrCoinProgramNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case errors.Is(err, errorz.ErrAccessDenied):
		return &echo.HTTPError{
			Code:    echo.ErrForbidden.Code,
			Message: err.Error(),
		}
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, coinProgram)
}
