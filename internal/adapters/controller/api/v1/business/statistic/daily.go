package statistic

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// TotalUsersDaily godoc
//
//	@Summary		Get business total users daily stats
//	@Description	Get business total users daily stats. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"period start date"
//	@Param			end_date	query		string	false	"period end date"
//	@Success		200			{object}	[]dto.BusinessStatsDailyTotalUsersResponse
//	@Failure		400			{object}	dto.HTTPStatus
//	@Failure		401			{object}	dto.HTTPStatus
//	@Failure		403			{object}	dto.HTTPStatus
//	@Failure		404			{object}	dto.HTTPStatus
//	@Failure		500			{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/stats/daily/total_users [get]
func (h *Handler) TotalUsersDaily(c echo.Context) error {
	var requestBody dto.BusinessStatsDailyTotalUsersRequest
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

	stats, err := h.statisticService.GetBusinessStatsDailyTotalUsers(
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

// ActiveUsersDaily godoc
//
//	@Summary		Get business active users daily stats
//	@Description	Get business active users daily stats. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"period start date"
//	@Param			end_date	query		string	false	"period end date"
//	@Success		200			{object}	[]dto.BusinessStatsDailyActiveUsersResponse
//	@Failure		400			{object}	dto.HTTPStatus
//	@Failure		401			{object}	dto.HTTPStatus
//	@Failure		403			{object}	dto.HTTPStatus
//	@Failure		404			{object}	dto.HTTPStatus
//	@Failure		500			{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/stats/daily/active_users [get]
func (h *Handler) ActiveUsersDaily(c echo.Context) error {
	var requestBody dto.BusinessStatsDailyActiveUsersRequest
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

	stats, err := h.statisticService.GetBusinessStatsDailyActiveUsers(
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
