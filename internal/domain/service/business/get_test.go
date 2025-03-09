package business_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/business"
	"testing"
)

func (suite *BusinessServiceTestSuite) TestGet() {
	t := suite.T()

	svc := business.NewService(nil, suite.db)

	cases := []struct {
		name       string
		businessID uuid.UUID

		wantResp *dto.BusinessReturn
		wantErr  error
	}{
		{
			name:       "success",
			businessID: suite.businessID,

			wantResp: &dto.BusinessReturn{
				ID:          suite.businessID,
				Name:        "Тест",
				Email:       "test@example.com",
				Description: "Тест",
			},
		},
		{
			name:       "business not found",
			businessID: uuid.New(),

			wantErr: errorz.ErrBusinessNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			resp, err := svc.Get(context.Background(), tCase.businessID)
			assert.Equal(t, tCase.wantErr, err)

			require.Equal(t, tCase.wantErr, err)
			require.Equal(t, tCase.wantResp, resp)
		})
	}
}
