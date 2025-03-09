package business_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/business"
	"testing"
)

func (suite *BusinessServiceTestSuite) TestRegister() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := business.NewService(nil, suite.db)

	cases := []struct {
		name string

		req *dto.BusinessRegisterRequest

		wantResp *dto.BusinessRegisterResponse
		wantErr  error
	}{
		{
			name: "success",

			req: &dto.BusinessRegisterRequest{
				Name:        "Test Business",
				Email:       "test2@example.com",
				Password:    "TestPass",
				Description: "Test Business",
			},
			wantResp: &dto.BusinessRegisterResponse{
				Name:        "Test Business",
				Email:       "test2@example.com",
				Description: "Test Business",
			},
		},
		{
			name: "email already exists",

			req: &dto.BusinessRegisterRequest{
				Name:        "Test Business",
				Email:       "test@example.com",
				Password:    "TestPass",
				Description: "Test Business",
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
