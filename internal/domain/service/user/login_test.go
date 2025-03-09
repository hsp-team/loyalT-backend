package user_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/user"
	"loyalit/internal/domain/service/user/mocks"
	"testing"
	"time"
)

func (suite *UserServiceTestSuite) TestLogin() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtConfigMock := mocks.NewMockJWTConfig(ctrl)

	svc := user.NewService(jwtConfigMock, suite.db)

	cases := []struct {
		name string

		req       *dto.UserLoginRequest
		mockSetup func()

		wantResp bool
		wantErr  error
	}{
		{
			name: "success",

			req: &dto.UserLoginRequest{
				Email:    "test@example.com",
				Password: "TestPass",
			},
			mockSetup: func() {
				jwtConfigMock.EXPECT().
					UserTokenSecret().
					Return("Test")
				jwtConfigMock.EXPECT().
					UserTokenExpiration().
					Return(time.Hour)
			},

			wantResp: true,
		},
		{
			name: "not found",

			req: &dto.UserLoginRequest{
				Email:    "notfound@example.com",
				Password: "TestPass",
			},
			mockSetup: func() {},

			wantErr: errorz.ErrUserNotFound,
		},
		{
			name: "password does not match",

			req: &dto.UserLoginRequest{
				Email:    "test@example.com",
				Password: "TestPass123",
			},
			mockSetup: func() {},

			wantErr: errorz.ErrPasswordDoesNotMatch,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.mockSetup()

			resp, err := svc.Login(context.Background(), tCase.req)
			assert.Equal(t, tCase.wantErr, err)
			assert.Equal(t, tCase.wantResp, resp != nil)
		})
	}
}
