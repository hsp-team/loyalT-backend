package coin_program_test

import (
	"bytes"
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

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinProgramServiceMock := mocks.NewMockcoinProgramService(ctrl)
	handler := coin_program.NewHandler(validator.New(), coinProgramServiceMock, nil)

	businessID := uuid.New()
	coinProgramParticipantID := uuid.New()

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
			body: &dto.CoinProgramCreateRequest{
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.CoinProgramCreateResponse{
						ID:        coinProgramParticipantID,
						Name:      "Тест",
						DayLimit:  1,
						CardColor: "#fff",
					}, nil)
			},

			expCode: http.StatusCreated,
			expResp: &dto.CoinProgramCreateResponse{
				ID:        coinProgramParticipantID,
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},
		},
		{
			name: "Empty name",
			body: &dto.CoinProgramCreateRequest{
				DayLimit:  1,
				CardColor: "#fff",
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
			name: "Empty day limit",
			body: &dto.CoinProgramCreateRequest{
				Name:      "Тест",
				CardColor: "#fff",
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[DayLimit]: '0' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Empty card color",
			body: &dto.CoinProgramCreateRequest{
				Name:     "Тест",
				DayLimit: 1,
			},
			businessID: businessID.String(),

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[CardColor]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			body: &dto.CoinProgramCreateRequest{
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Coin program already exists",
			body: &dto.CoinProgramCreateRequest{
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCoinProgramAlreadyExists)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "coin program already exists",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.CoinProgramCreateRequest{
				Name:      "Тест",
				DayLimit:  1,
				CardColor: "#fff",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				coinProgramServiceMock.EXPECT().
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

				var got *dto.CoinProgramCreateResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
