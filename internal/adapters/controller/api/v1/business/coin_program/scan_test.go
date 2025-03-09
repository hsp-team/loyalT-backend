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

func TestScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qrServiceMock := mocks.NewMockqrService(ctrl)
	handler := coin_program.NewHandler(validator.New(), nil, qrServiceMock)

	businessID := uuid.New()

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
			body: &dto.UserQRScanRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ScanUserQR(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.UserQRScanResponse{
						Username: "TESTUSER",
						Balance:  10,
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: &dto.UserQRScanResponse{
				Username: "TESTUSER",
				Balance:  10,
			},
		},
		{
			name:       "Invalid body",
			body:       &dto.UserQRScanRequest{},
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
			body: &dto.UserQRScanRequest{
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
			name: "User not found by QR",
			body: &dto.UserQRScanRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ScanUserQR(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrUserByQrNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "user by qr not found",
			},
			respErr: true,
		},
		{
			name: "Coin program not found",
			body: &dto.UserQRScanRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ScanUserQR(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrCoinProgramNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coin program not found",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.UserQRScanRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					ScanUserQR(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/scan", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Scan(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.UserQRScanResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}

func TestEnroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qrServiceMock := mocks.NewMockqrService(ctrl)
	handler := coin_program.NewHandler(validator.New(), nil, qrServiceMock)

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
			body: &dto.UserEnrollCoinRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					EnrollCoin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},

			expCode: http.StatusNoContent,
		},
		{
			name:       "Invalid body",
			body:       &dto.UserEnrollCoinRequest{},
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
			body: &dto.UserEnrollCoinRequest{
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
			name: "Coin program participant not found",
			body: &dto.UserEnrollCoinRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					EnrollCoin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errorz.ErrCoinProgramParticipantNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "coin program participant not found",
			},
			respErr: true,
		},
		{
			name: "User by QR not found",
			body: &dto.UserEnrollCoinRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					EnrollCoin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errorz.ErrUserByQrNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "user by qr not found",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.UserEnrollCoinRequest{
				Code: "TESTCODE",
			},
			businessID: businessID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					EnrollCoin(gomock.Any(), gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/enroll", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Enroll(c)

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
