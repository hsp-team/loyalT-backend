package business_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/business"
	"testing"
)

func (suite *BusinessServiceTestSuite) TestUpdate() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := business.NewService(nil, suite.db)

	cases := []struct {
		name       string
		businessID uuid.UUID
		req        *dto.BusinessUpdateRequest

		wantResp *dto.BusinessReturn
		wantErr  error
	}{
		{
			name:       "success",
			businessID: suite.businessID,
			req: &dto.BusinessUpdateRequest{
				Name:        "Updated Business",
				Description: "Updated business description",
			},

			wantResp: &dto.BusinessReturn{
				ID:          suite.businessID,
				Name:        "Updated Business",
				Email:       "test@example.com",
				Description: "Updated business description",
			},
		},
		{
			name:       "business not found",
			businessID: uuid.New(),
			req: &dto.BusinessUpdateRequest{
				Name:        "Updated Business",
				Description: "Updated business description",
			},

			wantErr: errorz.ErrBusinessNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			resp, err := svc.Update(context.Background(), tCase.req, tCase.businessID)
			assert.Equal(t, tCase.wantErr, err)
			assert.Equal(t, tCase.wantResp, resp)
		})
	}
}
