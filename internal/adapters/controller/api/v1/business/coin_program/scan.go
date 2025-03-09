package coin_program

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
)

// Scan godoc
//
//	@Summary		Scan user QR
//	@Description	Scan user QR to get user data and validate QR
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"QR code data"
//	@Success		200		{object}	dto.UserQRScanResponse
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program/scan/{code} [get]
func (h *Handler) Scan(c echo.Context) error {
	var requestBody dto.UserQRScanRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if validateErr := h.validator.ValidateData(requestBody); validateErr != nil {
		return validateErr
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	resp, err := h.qrService.ScanUserQR(
		c.Request().Context(),
		&requestBody,
		businessID,
	)
	switch {
	case errors.Is(err, errorz.ErrUserByQrNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case errors.Is(err, errorz.ErrCoinProgramNotFound):
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

	return c.JSON(http.StatusOK, resp)
}

// Enroll godoc
//
//	@Summary		Enroll coin to user by qr code
//	@Description	Add 1 coin to user coin program by qr code
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserEnrollCoinRequest	true	"User QR scan data"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		401		{object}	dto.HTTPStatus
//	@Failure		403		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business/coin_program/scan/enroll [post]
func (h *Handler) Enroll(c echo.Context) error {
	var requestBody dto.UserEnrollCoinRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if validateErr := h.validator.ValidateData(requestBody); validateErr != nil {
		return validateErr
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	err = h.qrService.EnrollCoin(
		c.Request().Context(),
		&requestBody,
		businessID,
	)
	switch {
	case errors.Is(err, errorz.ErrCoinProgramParticipantNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case errors.Is(err, errorz.ErrUserByQrNotFound):
		return &echo.HTTPError{
			Code:    echo.ErrNotFound.Code,
			Message: err.Error(),
		}
	case errors.Is(err, errorz.ErrUserScanLimitReached):
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

	return c.NoContent(204)
}
