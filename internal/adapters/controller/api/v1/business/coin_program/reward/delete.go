package reward

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Delete godoc
//
//	@Summary		Delete reward
//	@Description	Delete reward from coin program. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			reward_id	path		string	true	"Reward ID (uuid)"
//	@Success		204			{object}	nil
//	@Failure		400			{object}	dto.HTTPStatus
//	@Failure		401			{object}	dto.HTTPStatus
//	@Failure		403			{object}	dto.HTTPStatus
//	@Failure		404			{object}	dto.HTTPStatus
//	@Failure		500			{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program/reward/{reward_id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	var requestBody dto.RewardDeleteRequest
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
	err = h.rewardService.Delete(
		c.Request().Context(),
		&requestBody,
		businessID,
	)
	switch {
	case errors.Is(err, errorz.ErrRewardNotFound):
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

	return c.NoContent(204)
}
