package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Update godoc
//
//	@Summary		Update coin program
//	@Description	Update business coin program. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CoinProgramUpdateRequest	true	"Data for update"
//	@Success		200		{object}	dto.CoinProgramUpdateResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program [put]
func (h *Handler) Update(c echo.Context) error {
	var requestBody dto.CoinProgramUpdateRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if validateErr := h.validator.ValidateData(requestBody); validateErr != nil {
		return validateErr
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	coinProgram, err := h.coinProgramService.Update(c.Request().Context(), &requestBody, businessID)
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
