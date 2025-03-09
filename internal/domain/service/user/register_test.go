package user_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/user"
	"testing"
)

func (suite *UserServiceTestSuite) TestRegister() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := user.NewService(nil, suite.db)

	cases := []struct {
		name string

		req *dto.UserRegisterRequest

		wantResp *dto.UserRegisterResponse
		wantErr  error
	}{
		{
			name: "success",

			req: &dto.UserRegisterRequest{
				Name:     "Test User",
				Email:    "test2@example.com",
				Password: "TestPass",
			},
			wantResp: &dto.UserRegisterResponse{
				Name:  "Test User",
				Email: "test2@example.com",
			},
		},
		{
			name: "email already exists",

			req: &dto.UserRegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "TestPass",
			},

			wantErr: errorz.ErrEmailAlreadyExists,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			resp, err := svc.Register(context.Background(), tCase.req)
			assert.Equal(t, tCase.wantErr, err)

			assert.Equal(t, tCase.wantErr, err)

			if tCase.wantResp != nil {
				assert.Equal(t, tCase.wantResp.Name, resp.Name)
				assert.Equal(t, tCase.wantResp.Email, resp.Email)
			}
		})
	}
}
