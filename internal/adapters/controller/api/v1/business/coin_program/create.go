package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Create godoc
//
//	@Summary		Create new coin program
//	@Description	Create new coin program and return it. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CoinProgramCreateRequest	true	"Data for creation"
//	@Success		200		{object}	dto.CoinProgramCreateResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program [post]
func (h *Handler) Create(c echo.Context) error {
	var requestBody dto.CoinProgramCreateRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if errValidate := h.validator.ValidateData(requestBody); errValidate != nil {
		return errValidate
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	coinProgram, err := h.coinProgramService.Create(c.Request().Context(), &requestBody, businessID)
	switch {
	case errors.Is(err, errorz.ErrCoinProgramAlreadyExists):
		return &echo.HTTPError{
			Code:    echo.ErrConflict.Code,
			Message: errorz.ErrCoinProgramAlreadyExists.Error(),
		}
	case err != nil:
		return &echo.HTTPError{
			Code: echo.ErrInternalServerError.Code,

			Message: err.Error(),
		}
	}

	return c.JSON(201, coinProgram)
}
