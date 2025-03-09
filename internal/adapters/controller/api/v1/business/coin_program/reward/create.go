package reward

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Create godoc
//
//	@Summary		Create new reward
//	@Description	Create new reward for coin program. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Security		BusinessAuthCookie
//	@Param			request	body		dto.RewardCreateRequest	true	"Reward creation data"
//	@Success		201		{object}	dto.RewardCreateResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/business/coin_program/rewards [post]
func (h *Handler) Create(c echo.Context) error {
	var requestBody dto.RewardCreateRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrUnauthorized.Code,
		}
	}
	reward, err := h.rewardService.Create(
		c.Request().Context(),
		&requestBody,
		businessID,
	)
	switch {
	case errors.Is(err, errorz.ErrProgramNotFound):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrNotFound.Code,
		}
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	return c.JSON(201, reward)
}
