package reward_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program/reward"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program/reward/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := reward.NewHandler(validator.New(), rewardServiceMock)

	userID := uuid.New()
	coinProgramParticipantID := uuid.New()
	rewardID1 := uuid.New()
	rewardID2 := uuid.New()

	cases := []struct {
		name   string
		body   interface{}
		userID string

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name: "success",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserList(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]dto.RewardReturn{
						{
							ID:       rewardID1,
							Name:     "Тест1",
							Cost:     10,
							ImageURL: "Картинка1",
						},
						{
							ID:       rewardID2,
							Name:     "Тест2",
							Cost:     5,
							ImageURL: "Картинка2",
						},
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: []dto.RewardReturn{
				{
					ID:       rewardID1,
					Name:     "Тест1",
					Cost:     10,
					ImageURL: "Картинка1",
				},
				{
					ID:       rewardID2,
					Name:     "Тест2",
					Cost:     5,
					ImageURL: "Картинка2",
				},
			},
		},
		{
			name: "Invalid limit and offset",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  -1,
				Offset: -1,
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Limit]: '-1' | Needs to implement 'gte' and [Offset]: '-1' | Needs to implement 'gte'",
			},
			respErr: true,
		},
		{
			name: "Invalid limit and offset",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  5,
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
			name: "coin program participant not found",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserList(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCoinProgramParticipantNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coin program participant not found",
			},
			respErr: true,
		},
		{
			name: "access denied",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserList(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrAccessDenied)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "access denied",
			},
			respErr: true,
		},
		{
			name: "internal error",
			body: dto.RewardUserListRequest{
				ID:     coinProgramParticipantID,
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserList(gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Set("user_id", tCase.userID)

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
