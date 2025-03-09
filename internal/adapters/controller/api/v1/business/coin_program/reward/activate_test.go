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

func TestActivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qrServiceMock := mocks.NewMockqrService(ctrl)
	handler := reward.NewHandler(validator.New(), nil, qrServiceMock)

	businessID := uuid.New()
	rewardID := uuid.New()

	cases := []struct {
		name       string
		body       interface{}
		businessID string

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name: "Success",
			body: &dto.RewardActivateRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ActivateUserReward(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.RewardReturn{
						ID:       rewardID,
						Name:     "Тест",
						Cost:     10,
						ImageURL: "https://image.jpg",
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: &dto.RewardReturn{
				ID:       rewardID,
				Name:     "Тест",
				Cost:     10,
				ImageURL: "https://image.jpg",
			},
		},
		{
			name:       "Invalid body",
			body:       &dto.RewardActivateRequest{},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Code]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			body: &dto.RewardActivateRequest{
				Code: "TESTCODE",
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Coupon not found",
			body: &dto.RewardActivateRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ActivateUserReward(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCouponNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coupon not found",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.RewardActivateRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ActivateUserReward(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
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

			req := httptest.NewRequest(http.MethodPost, "/activate", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Activate(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.RewardReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
