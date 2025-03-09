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
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := reward.NewHandler(validator.New(), rewardServiceMock, nil)

	businessID := uuid.New()
	rewardID1 := uuid.New()
	rewardID2 := uuid.New()

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
			body: &dto.RewardListRequest{
				Limit:  10,
				Offset: 0,
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]dto.RewardReturn{
						{
							ID:       rewardID1,
							Name:     "Test Reward 1",
							Cost:     100,
							ImageURL: "https://image.url/1",
						},
						{
							ID:       rewardID2,
							Name:     "Test Reward 2",
							Cost:     150,
							ImageURL: "https://image.url/2",
						},
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: []dto.RewardReturn{
				{
					ID:       rewardID1,
					Name:     "Test Reward 1",
					Cost:     100,
					ImageURL: "https://image.url/1",
				},
				{
					ID:       rewardID2,
					Name:     "Test Reward 2",
					Cost:     150,
					ImageURL: "https://image.url/2",
				},
			},
		},
		{
			name: "Invalid limit and offset",
			body: &dto.RewardListRequest{
				Limit:  -1,
				Offset: -1,
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Limit]: '-1' | Needs to implement 'gte' and [Offset]: '-1' | Needs to implement 'gte'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			body: &dto.RewardListRequest{
				Limit:  10,
				Offset: 0,
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Internal error in service",
			body: &dto.RewardListRequest{
				Limit:  10,
				Offset: 0,
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodGet, "/list", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.List(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got []dto.RewardReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
