package reward

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Activate godoc
//
//	@Summary		Activate scanned coupon code
//	@Description	Activate scanned coupon code. Requires business_auth_token cookie for authentication.
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Security		BusinessAuthCookie
//	@Param			request	body		dto.RewardActivateRequest	true	"Reward activation data"
//	@Success		200		{object}	dto.RewardReturn
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Router			/business/coin_program/rewards/activate [post]
func (h *Handler) Activate(c echo.Context) error {
	var requestBody dto.RewardActivateRequest
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

	reward, err := h.qrService.ActivateUserReward(c.Request().Context(), &requestBody, businessID)
	switch {
	case errors.Is(err, errorz.ErrCouponNotFound):
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

	return c.JSON(http.StatusOK, reward)
}
