package reward

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Buy godoc
//
//	@Summary		Buy reward by its id
//	@Description	Buy reward by its id. Requires user_auth_token cookie for authentication.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		UserAuthCookie
//	@Param			coin_program_participant_id	path		string	true	"Coin program ID (uuid)"
//	@Param			reward_id					path		string	true	"Reward ID (uuid)"
//	@Success		200							{object}	dto.RewardBuyResponse
//	@Failure		400							{object}	dto.HTTPStatus
//	@Failure		401							{object}	dto.HTTPStatus
//	@Failure		402							{object}	dto.HTTPStatus
//	@Failure		403							{object}	dto.HTTPStatus
//	@Failure		404							{object}	dto.HTTPStatus
//	@Failure		500							{object}	dto.HTTPStatus
//	@Router			/user/coin_programs/{coin_program_id}/rewards/{reward_id}/buy [post]
func (h *Handler) Buy(c echo.Context) error {
	var requestBody dto.RewardBuyRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrUnauthorized.Code,
		}
	}

	code, err := h.rewardService.Buy(c.Request().Context(), &requestBody, userID)
	switch {
	case errors.Is(err, errorz.ErrCoinProgramParticipantNotFound):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrNotFound.Code,
		}
	case errors.Is(err, errorz.ErrAccessDenied):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrForbidden.Code,
		}
	case errors.Is(err, errorz.ErrNotEnoughCoins):
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrPaymentRequired.Code,
		}
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	return c.JSON(http.StatusOK, code)
}
