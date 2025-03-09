package reward

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/entity/dto"
)

// List godoc
//
//	@Summary		List rewards
//	@Description	List rewards for coin program. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int					false	"Number of items to return"	minimum(1)	default(10)
//	@Param			offset	query		int					false	"Number of items to skip"	minimum(0)	default(0)
//	@Success		200		{object}	[]dto.RewardReturn	"List of rewards"
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program/rewards [get]
func (h *Handler) List(c echo.Context) error {
	var requestBody dto.RewardListRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	if requestBody.Limit == 0 {
		requestBody.Limit = 10
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrUnauthorized.Code,
		}
	}
	rewards, err := h.rewardService.List(
		c.Request().Context(),
		&requestBody,
		businessID,
	)
	switch {
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	return c.JSON(200, rewards)
}
