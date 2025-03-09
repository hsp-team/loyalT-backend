package user_test

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"loyalit/internal/adapters/repository/postgres/ent"
	entuser "loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/test/testhelper"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	pgContainer *testhelper.PostgresContainer
	db          *ent.Client
	ctx         context.Context
	userID      uuid.UUID
}

func (suite *UserServiceTestSuite) SetupSuite() {
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

func (suite *UserServiceTestSuite) TearDownSuite() {
	if err := suite.db.Close(); err != nil {
		panic(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		panic(err)
	}
}

func (suite *UserServiceTestSuite) SetupTest() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("TestPass"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user, err := suite.db.User.Create().
		SetName("Тест").
		SetEmail("test@example.com").
		SetPassword(string(hashedPassword)).
		Save(suite.ctx)
	if err != nil {
		panic(err)
	}
	suite.userID = user.ID
}

func (suite *UserServiceTestSuite) TearDownTest() {
	_, err := suite.db.User.Delete().Where(entuser.ID(suite.userID)).Exec(suite.ctx)
	if err != nil {
		panic(err)
	}
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
