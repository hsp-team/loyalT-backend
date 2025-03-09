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

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := reward.NewHandler(validator.New(), rewardServiceMock, nil)

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
			body: &dto.RewardCreateRequest{
				Name:     "Test Reward",
				Cost:     100,
				ImageURL: "https://image.url",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.RewardCreateResponse{
						ID:       rewardID,
						Name:     "Test Reward",
						Cost:     100,
						ImageURL: "https://image.url",
					}, nil)
			},

			expCode: http.StatusCreated,
			expResp: &dto.RewardCreateResponse{
				ID:       rewardID,
				Name:     "Test Reward",
				Cost:     100,
				ImageURL: "https://image.url",
			},
		},
		{
			name: "Empty Name",
			body: &dto.RewardCreateRequest{
				Cost:     100,
				ImageURL: "https://image.url",
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Name]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Empty Cost",
			body: &dto.RewardCreateRequest{
				Name:     "Тест",
				ImageURL: "https://image.url",
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Cost]: '0' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Empty ImageURL",
			body: &dto.RewardCreateRequest{
				Name: "Тест",
				Cost: 100,
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[ImageURL]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid ImageURL",
			body: &dto.RewardCreateRequest{
				Name:     "Тест",
				Cost:     100,
				ImageURL: "invalid",
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[ImageURL]: 'invalid' | Needs to implement 'url'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			body: &dto.RewardCreateRequest{
				Name:     "Test Reward",
				Cost:     100,
				ImageURL: "https://image.url",
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Program not found",
			body: &dto.RewardCreateRequest{
				Name:     "Test Reward",
				Cost:     100,
				ImageURL: "https://image.url",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrProgramNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "program not found",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.RewardCreateRequest{
				Name:     "Test Reward",
				Cost:     100,
				ImageURL: "https://image.url",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Create(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.RewardCreateResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
