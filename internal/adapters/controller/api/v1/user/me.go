package user

import (
	"errors"
	"loyalit/internal/domain/common/errorz"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Me godoc
//
//	@Summary		Get user profile
//	@Description	Get user profile
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserGetResponse
//	@Failure		400	{object}	dto.HTTPStatus
//	@Failure		401	{object}	dto.HTTPStatus
//	@Failure		404	{object}	dto.HTTPStatus
//	@Failure		500	{object}	dto.HTTPStatus
//	@Router			/user/me [get]
func (h *Handler) Me(c echo.Context) error {
	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	ctx := c.Request().Context()
	user, err := h.userService.Get(ctx, userID)
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

	// Get user statistics before returning response
	userStats, err := h.statisticService.GetUserStatistics(ctx, userID)
	if err != nil {
		// Log the error but don't fail the request
		c.Logger().Errorf("Failed to get user statistics: %v", err)
	} else {
		user.Statistics = *userStats
	}

	return c.JSON(http.StatusOK, user)
}
