package reward

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// List godoc
//
//	@Summary		List rewards for user
//	@Description	List rewards for user by coin program id. Requires business_auth_token cookie for authentication.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			coin_program_participant_id	path		string				true	"Coin program participant id"
//	@Param			limit						query		int					false	"Number of items to return"	minimum(1)	default(10)
//	@Param			offset						query		int					false	"Number of items to skip"	minimum(0)	default(0)
//	@Success		200							{object}	[]dto.RewardReturn	"List of rewards"
//	@Failure		400							{object}	dto.HTTPStatus
//	@Failure		401							{object}	dto.HTTPStatus
//	@Failure		500							{object}	dto.HTTPStatus
//	@Security		UserAuthCookie
//	@Router			/user/coin_programs/{coin_program_participant_id}/rewards [get]
func (h *Handler) List(c echo.Context) error {
	var requestBody dto.RewardUserListRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	if requestBody.Limit == 0 {
		requestBody.Limit = 10
	}

	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrUnauthorized.Code,
		}
	}
	rewards, err := h.rewardService.UserList(
		c.Request().Context(),
		&requestBody,
		userID,
	)
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
	case err != nil:
		return &echo.HTTPError{
			Message: err.Error(),
			Code:    echo.ErrInternalServerError.Code,
		}
	}

	return c.JSON(200, rewards)
}
