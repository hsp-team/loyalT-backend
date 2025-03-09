package business

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

// Update godoc
//
//	@Summary		Update business
//	@Description	Update business
//	@Tags			Business
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.BusinessUpdateRequest	true	"Business data for update"
//	@Success		201		{object}	dto.BusinessReturn
//	@Failure		400		{object}	dto.HTTPStatus
//	@Failure		404		{object}	dto.HTTPStatus
//	@Failure		500		{object}	dto.HTTPStatus
//	@Security		BusinessAuthCookie
//	@Router			/business [put]
func (h *Handler) Update(c echo.Context) error {
	var requestBody dto.BusinessUpdateRequest
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := h.validator.ValidateData(requestBody); err != nil {
		return err
	}

	businessID, err := uuid.Parse(c.Get("business_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: err.Error(),
		}
	}

	business, err := h.businessService.Update(c.Request().Context(), &requestBody, businessID)
	switch {
	case errors.Is(err, errorz.ErrBusinessNotFound):
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

	return c.JSON(200, business)
}
