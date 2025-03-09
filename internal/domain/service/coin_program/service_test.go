package coin_program_test

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"loyalit/internal/adapters/repository/postgres/ent"
	entbusiness "loyalit/internal/adapters/repository/postgres/ent/business"
	entcoinprogram "loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/test/testhelper"
	"testing"
)

type CoinProgramServiceTestSuite struct {
	suite.Suite
	pgContainer   *testhelper.PostgresContainer
	db            *ent.Client
	ctx           context.Context
	coinProgramID uuid.UUID
	businessID    uuid.UUID
}

func (suite *CoinProgramServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelper.CreatePostgresContainer(suite.ctx)
	if err != nil {
		panic(err)
	}
	suite.pgContainer = pgContainer

	db, err := ent.Open(dialect.Postgres, suite.pgContainer.ConnStr)
	if err != nil {
		panic(err)
	}

	suite.db = db

	if errMigrate := suite.db.Schema.Create(
		context.Background(),
	); errMigrate != nil {
		panic(err)
	}

}

func (suite *CoinProgramServiceTestSuite) TearDownSuite() {
	if err := suite.db.Close(); err != nil {
		panic(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		panic(err)
	}
}

func (suite *CoinProgramServiceTestSuite) SetupTest() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("TestPass"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	business, err := suite.db.Business.Create().
		SetName("Тест").
		SetEmail("test@example.com").
		SetPassword(string(hashedPassword)).
		SetDescription("Тест").
		Save(suite.ctx)
	if err != nil {
		panic(err)
	}

	coinProgram, err := suite.db.CoinProgram.
		Create().
		SetName("Тест").
		SetDescription("Тест").
		SetDayLimit(1).
		SetCardColor("#fff").
		SetBusinessID(business.ID).
		Save(suite.ctx)
	if err != nil {
		panic(err)
	}
	suite.businessID = business.ID
	suite.coinProgramID = coinProgram.ID
}

func (suite *CoinProgramServiceTestSuite) TearDownTest() {
	_, err := suite.db.CoinProgram.Delete().Where(entcoinprogram.ID(suite.coinProgramID)).Exec(suite.ctx)
	if err != nil {
		panic(err)
	}
	_, err = suite.db.Business.Delete().Where(entbusiness.ID(suite.businessID)).Exec(suite.ctx)
	if err != nil {
		panic(err)
	}
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CoinProgramServiceTestSuite))
}
