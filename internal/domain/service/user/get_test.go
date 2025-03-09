package user_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/user"
	"testing"
)

func (suite *UserServiceTestSuite) TestGet() {
	t := suite.T()

	svc := user.NewService(nil, suite.db)

	cases := []struct {
		name string

		userID uuid.UUID

		wantResp *dto.UserGetResponse
		wantErr  error
	}{
		{
			name: "success",

			userID: suite.userID,

			wantResp: &dto.UserGetResponse{
				ID:    suite.userID,
				Name:  "Тест",
				Email: "test@example.com",
			},
		},
		{
			name: "not found",

			userID: uuid.Nil,

			wantErr: errorz.ErrUserNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			resp, err := svc.Get(context.Background(), tCase.userID)
			assert.Equal(t, tCase.wantErr, err)
			assert.Equal(t, tCase.wantResp, resp)
		})
	}
}
