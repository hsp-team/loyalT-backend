package coin_program_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/service/coin_program"
	"testing"
)

func (suite *CoinProgramServiceTestSuite) TestCreate() {
	t := suite.T()

	svc := coin_program.NewService(suite.db)

	cases := []struct {
		name string

		businessID uuid.UUID
		req        *dto.CoinProgramCreateRequest

		wantErr error
	}{
		{
			name: "success",

			businessID: suite.businessID,
			req: &dto.CoinProgramCreateRequest{
				Name:        "Тест",
				Description: "Тест",
				DayLimit:    1,
				CardColor:   "#fff",
			},

			wantErr: errorz.ErrCoinProgramAlreadyExists,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tCase.req, tCase.businessID)
			assert.Equal(t, tCase.wantErr, err)

			assert.Equal(t, tCase.wantErr, err)
		})
	}
}
