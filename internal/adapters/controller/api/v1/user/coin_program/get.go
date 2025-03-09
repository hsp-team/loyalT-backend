package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Get godoc
//
//	@Summary		Get user's coin program by id
//	@Description	Get user's coin program by ID. Requires business_auth_token cookie for authentication.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			coin_program_participant_id	path		string	true	"Coin program ID (uuid)"
//	@Success		200							{object}	dto.CoinProgramParticipantReturn
//	@Failure		400							{object}	dto.HTTPStatus
//	@Failure		401							{object}	dto.HTTPStatus
//	@Failure		403							{object}	dto.HTTPStatus
//	@Failure		404							{object}	dto.HTTPStatus
//	@Failure		500							{object}	dto.HTTPStatus
//	@Security		UserAuthCookie
//	@Router			/user/coin_programs/{coin_program_participant_id} [get]
func (h *Handler) Get(c echo.Context) error {
	var requestBody dto.CoinProgramParticipantGetRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if validateErr := h.validator.ValidateData(requestBody); validateErr != nil {
		return validateErr
	}

	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	coinProgram, err := h.coinProgramParticipantService.Get(c.Request().Context(), requestBody.ID, userID)
	switch {
	case errors.Is(err, errorz.ErrCoinProgramNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case errors.Is(err, errorz.ErrAccessDenied):
		return &echo.HTTPError{
			Code:    echo.ErrForbidden.Code,
			Message: err.Error(),
		}
	case err != nil:
		return &echo.HTTPError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, coinProgram)
}
