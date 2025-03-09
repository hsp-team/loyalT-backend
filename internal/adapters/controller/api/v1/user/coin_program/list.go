package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// List godoc
//
//	@Summary		List user's coin programs
//	@Description	Get user's coin programs. Requires business_auth_token cookie for authentication.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	true	"Limit"
//	@Param			offset	query		int	true	"Offset"
//	@Success		200		{object}	[]dto.CoinProgramParticipantReturn
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		UserAuthCookie
//	@Router			/user/coin_programs [get]
func (h *Handler) List(c echo.Context) error {
	var requestBody dto.CoinProgramParticipantListRequest
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

	coinProgramParticipants, err := h.coinProgramParticipantService.UserList(c.Request().Context(), requestBody, userID)
	switch {
	case errors.Is(err, errorz.ErrUserNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, coinProgramParticipants)
}
