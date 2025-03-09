package coin_program

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// ListAvailable godoc
//
//	@Summary		Get rewards available for user
//	@Description	Get rewards available for user. Requires business_auth_token cookie for authentication.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	true	"Limit"
//	@Param			offset	query		int	true	"Offset"
//	@Success		200		{object}	[]dto.CoinProgramWithRewardsReturn
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		UserAuthCookie
//	@Router			/user/coin_programs/rewards/available [get]
func (h *Handler) ListAvailable(c echo.Context) error {
	var requestBody dto.CoinProgramParticipantListAvailableRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if validateErr := h.validator.ValidateData(requestBody); validateErr != nil {
		return validateErr
	}

	if requestBody.Limit == 0 {
		requestBody.Limit = 10
	}

	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	rewards, err := h.rewardService.UserListAvailable(c.Request().Context(), &requestBody, userID)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, rewards)
}
