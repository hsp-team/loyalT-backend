package business_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/business/mocks"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/business"
	"testing"
	"time"
)

func (suite *BusinessServiceTestSuite) TestLogin() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtConfigMock := mocks.NewMockJWTConfig(ctrl)

	svc := business.NewService(jwtConfigMock, suite.db)

	cases := []struct {
		name string

		req       *dto.BusinessLoginRequest
		mockSetup func()

		wantResp bool
		wantErr  error
	}{
		{
			name: "success",

			req: &dto.BusinessLoginRequest{
				Email:    "test@example.com",
				Password: "TestPass",
			},
			mockSetup: func() {
				jwtConfigMock.EXPECT().
					BusinessTokenSecret().
					Return("Test")
				jwtConfigMock.EXPECT().
					BusinessTokenExpiration().
					Return(time.Hour)
			},

			wantResp: true,
		},
		{
			name: "not found",

			req: &dto.BusinessLoginRequest{
				Email:    "notfound@example.com",
				Password: "TestPass",
			},
			mockSetup: func() {},

			wantErr: errorz.ErrBusinessNotFound,
		},
		{
			name: "password does not match",

			req: &dto.BusinessLoginRequest{
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
