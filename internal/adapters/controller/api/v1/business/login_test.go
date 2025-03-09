package business_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/business"
	"loyalit/internal/adapters/controller/api/v1/business/mocks"
	"loyalit/internal/domain/common/errorz"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	businessServiceMock := mocks.NewMockbusinessService(ctrl)
	handler := business.NewHandler(nil, validator.New(), businessServiceMock, true)

	cases := []struct {
		name string
		body interface{}

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError // ожидаемая ошибка
		respErr bool            // флаг, указывающий на ошибку
	}{
		{
			name: "Success",
			body: dto.BusinessLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(&dto.BusinessLoginResponse{
						Token: "access-token",
					}, nil)
			},
			expCode: http.StatusOK,
			expResp: &dto.HTTPStatus{
				Code:    http.StatusOK,
				Message: "Successfully logged in as business",
			},
			respErr: false,
		},
		{
			name: "Invalid email",
			body: dto.BusinessLoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Email]: 'invalid-email' | Needs to implement 'email'",
			},
			respErr: true,
		},
		{
			name: "Empty password",
			body: dto.BusinessLoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Password]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid credentials",
			body: dto.BusinessLoginRequest{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrPasswordDoesNotMatch)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "password does not match",
			},
			respErr: true,
		},
		{
			name: "User not found",
			body: dto.BusinessLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrBusinessNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "business not found",
			},
			respErr: true,
		},
		{
			name: "Internal server error",
			body: dto.BusinessLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			err = handler.Login(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.HTTPStatus
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
