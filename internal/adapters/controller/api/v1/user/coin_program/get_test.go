package coin_program_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinProgramParticipantMock := mocks.NewMockcoinProgramParticipantService(ctrl)
	handler := coin_program.NewHandler(validator.New(), coinProgramParticipantMock, nil)

	userID := uuid.New()
	coinProgramParticipantID := uuid.New()

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
			name: "Success",
			body: &dto.CoinProgramParticipantGetRequest{
				ID: coinProgramParticipantID,
			},
			userID: userID.String(),

			mockSetup: func() {
				coinProgramParticipantMock.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.CoinProgramParticipantReturn{
						ID:        coinProgramParticipantID,
						Name:      "Тест",
						Balance:   10,
						CardColor: "#fff",
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: &dto.CoinProgramParticipantReturn{
				ID:        coinProgramParticipantID,
				Name:      "Тест",
				Balance:   10,
				CardColor: "#fff",
			},
		},
		{
			name: "Invalid coin program id",

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[ID]: '00000000-0000-0000-0000-000000000000' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid coin program id",
			body: &dto.CoinProgramParticipantGetRequest{
				ID: coinProgramParticipantID,
			},

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "User not found",
			body: &dto.CoinProgramParticipantGetRequest{
				ID: coinProgramParticipantID,
			},
			userID: userID.String(),

			mockSetup: func() {
				coinProgramParticipantMock.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCoinProgramNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coin program not found",
			},
			respErr: true,
		},
		{
			name: "Access denied",
			body: &dto.CoinProgramParticipantGetRequest{
				ID: coinProgramParticipantID,
			},
			userID: userID.String(),

			mockSetup: func() {
				coinProgramParticipantMock.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrAccessDenied)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "access denied",
			},
			respErr: true,
		},
		{
			name: "Internal server error",
			body: &dto.CoinProgramParticipantGetRequest{
				ID: coinProgramParticipantID,
			},
			userID: userID.String(),

			mockSetup: func() {
				coinProgramParticipantMock.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodGet, "/get", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("user_id", tCase.userID)

			err = handler.Get(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.CoinProgramParticipantReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
