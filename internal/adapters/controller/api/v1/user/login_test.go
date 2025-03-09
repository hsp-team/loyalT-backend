package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/user"
	"loyalit/internal/adapters/controller/api/v1/user/mocks"
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

	userServiceMock := mocks.NewMockuserService(ctrl)
	handler := user.NewHandler(nil, validator.New(), userServiceMock, nil, nil, true)

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
			body: dto.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				userServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(&dto.UserLoginResponse{
						Token: "access-token",
					}, nil)
			},
			expCode: http.StatusOK,
			expResp: &dto.HTTPStatus{
				Code:    http.StatusOK,
				Message: "Successfully logged in as user",
			},
			respErr: false,
		},
		{
			name: "Invalid email",
			body: dto.UserLoginRequest{
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
			body: dto.UserLoginRequest{
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
			body: dto.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockSetup: func() {
				userServiceMock.EXPECT().
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
			body: dto.UserLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				userServiceMock.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrUserNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "user not found",
			},
			respErr: true,
		},
		{
			name: "Internal server error",
			body: dto.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				userServiceMock.EXPECT().
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
