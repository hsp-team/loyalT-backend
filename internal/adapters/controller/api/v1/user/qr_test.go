package user_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/user"
	"loyalit/internal/adapters/controller/api/v1/user/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qrServiceMock := mocks.NewMockqrService(ctrl)
	handler := user.NewHandler(nil, validator.New(), nil, qrServiceMock, nil, true)

	userID := uuid.New()

	cases := []struct {
		name   string
		userID string

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name:   "Success",
			userID: userID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					GetUserQR(gomock.Any(), gomock.Any()).
					Return(&dto.QRGetResponse{
						Data: "TESTDATA",
					}, nil)
			},
			expCode: http.StatusOK,
			expResp: &dto.QRGetResponse{
				Data: "TESTDATA",
			},
		},
		{
			name: "Invalid user id",

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name:   "User not found",
			userID: userID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					GetUserQR(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrUserNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "user not found",
			},
			respErr: true,
		},
		{
			name:   "Internal server error",
			userID: userID.String(),

			mockSetup: func() {
				qrServiceMock.EXPECT().
					GetUserQR(gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodGet, "/qr", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("user_id", tCase.userID)

			err := handler.GetQR(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.QRGetResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
