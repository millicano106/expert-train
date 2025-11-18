package customer

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	pgContainer *postgres.PostgresContainer
	repository  *Repository
	ctx         context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := postgres.RunContainer(suite.ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	connStr, err := pgContainer.ConnectionString(suite.ctx, "sslmode=disable")
	repository, err := NewRepository(suite.ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository
}

func (suite *CustomerRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *CustomerRepoTestSuite) TestCreateCustomer() {
	t := suite.T()

	customer, err := suite.repository.CreateCustomer(suite.ctx, Customer{
		Name:  "Henry",
		Email: "henry@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, customer.Id)
}

func (suite *CustomerRepoTestSuite) TestGetCustomerByEmail() {
	t := suite.T()

	customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "john@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "John", customer.Name)
	assert.Equal(t, "john@gmail.com", customer.Email)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}
