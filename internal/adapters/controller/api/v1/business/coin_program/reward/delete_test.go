package reward_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/business/coin_program/reward"
	"loyalit/internal/adapters/controller/api/v1/business/coin_program/reward/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := reward.NewHandler(validator.New(), rewardServiceMock, nil)

	businessID := uuid.New()

	cases := []struct {
		name       string
		body       interface{}
		businessID string

		mockSetup func()

		expCode int
		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name: "Success",
			body: &dto.RewardDeleteRequest{
				ID: uuid.New(),
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},

			expCode: http.StatusNoContent,
		},
		{
			name:       "Invalid body",
			body:       &dto.RewardDeleteRequest{},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[ID]: '00000000-0000-0000-0000-000000000000' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			body: &dto.RewardDeleteRequest{
				ID: uuid.New(),
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Reward not found",
			body: &dto.RewardDeleteRequest{
				ID: uuid.New(),
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errorz.ErrRewardNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "reward not found",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.RewardDeleteRequest{
				ID: uuid.New(),
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("internal error"))
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "internal error",
			},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.mockSetup()

			bodyBytes, err := json.Marshal(tCase.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Delete(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)
			}
		})
	}
}
