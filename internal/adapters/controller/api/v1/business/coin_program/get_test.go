package coin_program_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/business/coin_program"
	"loyalit/internal/adapters/controller/api/v1/business/coin_program/mocks"
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

	coinProgramServiceMock := mocks.NewMockcoinProgramService(ctrl)
	handler := coin_program.NewHandler(validator.New(), coinProgramServiceMock, nil)

	businessID := uuid.New()
	coinProgramID := uuid.New()

	cases := []struct {
		name       string
		businessID string

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name:       "Success",
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(&dto.CoinProgramReturn{
						ID:        coinProgramID,
						Name:      "Тест",
						DayLimit:  1,
						CardColor: "#fff",
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: &dto.CoinProgramReturn{
				ID:        coinProgramID,
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},
		},
		{
			name: "Invalid business id",

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name:       "Coin program not found",
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCoinProgramNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coin program not found",
			},
			respErr: true,
		},
		{
			name:       "Access denied",
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrAccessDenied)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "access denied",
			},
			respErr: true,
		},
		{
			name:       "Internal error",
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Get(gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodGet, "/get", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err := handler.Get(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.CoinProgramReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
