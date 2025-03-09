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

func (suite *CoinProgramServiceTestSuite) TestGet() {
	t := suite.T()

	svc := coin_program.NewService(suite.db)

	cases := []struct {
		name string

		businessID uuid.UUID

		wantResp *dto.CoinProgramReturn
		wantErr  error
	}{
		{
			name: "success",

			businessID: suite.businessID,

			wantResp: &dto.CoinProgramReturn{
				ID:          suite.coinProgramID,
				Name:        "Тест",
				Description: "Тест",
				DayLimit:    1,
				CardColor:   "#fff",
			},
		},
		{
			name: "not found",

			businessID: uuid.New(),

			wantErr: errorz.ErrCoinProgramNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			resp, err := svc.Get(context.Background(), tCase.businessID)
			assert.Equal(t, tCase.wantErr, err)
			assert.Equal(t, tCase.wantResp, resp)
		})
	}
}
