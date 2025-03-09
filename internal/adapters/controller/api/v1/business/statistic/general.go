package statistic

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Business godoc
//
//	@Summary		Get business stats
//	@Description	Get business general stats. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"period start date"
//	@Param			end_date	query		string	false	"period end date"
//	@Success		200			{object}	dto.BusinessStatsResponse
//	@Failure		400			{object}	dto.HTTPStatus
//	@Failure		401			{object}	dto.HTTPStatus
//	@Failure		403			{object}	dto.HTTPStatus
//	@Failure		404			{object}	dto.HTTPStatus
//	@Failure		500			{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/stats/general [get]
func (h *Handler) Business(c echo.Context) error {
	var requestBody dto.BusinessStatsRequest
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

	stats, err := h.statisticService.GetBusinessStats(
		c.Request().Context(),
		businessID,
		&requestBody,
	)
	switch {
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, stats)
}

// CoinProgram godoc
//
//	@Summary		Get business coin program stats
//	@Description	Get business coin program stats. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"period start date"
//	@Param			end_date	query		string	false	"period end date"
//	@Success		200			{object}	dto.BusinessCoinProgramStatsResponse
//	@Failure		400			{object}	dto.HTTPStatus
//	@Failure		401			{object}	dto.HTTPStatus
//	@Failure		403			{object}	dto.HTTPStatus
//	@Failure		404			{object}	dto.HTTPStatus
//	@Failure		500			{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/stats/coin_program [get]
func (h *Handler) CoinProgram(c echo.Context) error {
	var requestBody dto.BusinessCoinProgramStatsRequest
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

	stats, err := h.statisticService.GetBusinessCoinProgramStats(
		c.Request().Context(),
		businessID,
		&requestBody,
	)
	switch {
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, stats)
}
