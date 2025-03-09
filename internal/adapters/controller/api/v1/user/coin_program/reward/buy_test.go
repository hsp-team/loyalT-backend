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

func TestBuy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := reward.NewHandler(validator.New(), rewardServiceMock)

	userID := uuid.New()
	coinProgramParticipantID := uuid.New()
	rewardID := uuid.New()

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
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Buy(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.RewardBuyResponse{
						Code: "TESTCODE",
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: &dto.RewardBuyResponse{
				Code: "TESTCODE",
			},
		},
		{
			name: "empty coin program participant id",
			body: dto.RewardBuyRequest{
				ID: rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[CoinProgramParticipantID]: '00000000-0000-0000-0000-000000000000' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "empty reward id",
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
			},
			userID: userID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[ID]: '00000000-0000-0000-0000-000000000000' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "invalid user id",
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
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
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Buy(gomock.Any(), gomock.Any(), gomock.Any()).
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
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Buy(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrAccessDenied)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "access denied",
			},
			respErr: true,
		},
		{
			name: "not enough coins",
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Buy(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrNotEnoughCoins)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusPaymentRequired,
				Message: "not enough coins",
			},
			respErr: true,
		},
		{
			name: "internal error",
			body: dto.RewardBuyRequest{
				CoinProgramParticipantID: coinProgramParticipantID,
				ID:                       rewardID,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					Buy(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/buy", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("user_id", tCase.userID)

			err = handler.Buy(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.RewardBuyResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
