package config

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DatabaseTestSuite struct {
	suite.Suite
	ctx                context.Context
	pgContainer        testcontainers.Container
	pgConnectionString string
	postgresClient     PostgresClient
	config             Config
	migrationURL       string
}

func (suite *DatabaseTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	containerPort := "5432"
	req := testcontainers.ContainerRequest{
		Image: "postgres:16-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "Numeris_DB",
		},
		ExposedPorts: []string{containerPort + "/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort(nat.Port(containerPort)),
		).WithDeadline(5 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.NoError(err)

	host, err := container.Host(suite.ctx)
	suite.NoError(err)

	mappedPort, err := container.MappedPort(suite.ctx, nat.Port(containerPort))
	suite.NoError(err)

	connStr := fmt.Sprintf("postgresql://postgres:postgres@%s:%s/Numeris_DB?sslmode=disable", host, mappedPort.Port())

	suite.pgContainer = container
	suite.pgConnectionString = connStr
	suite.migrationURL = "file://../../migrations"
	suite.postgresClient = PostgresClient{
		DBSource: connStr,
	}
	suite.config = Load("dev", "0.0.0.0:3000", connStr)
}

func (suite *DatabaseTestSuite) TearDownSuite() {
	if suite.pgContainer != nil {
		err := suite.pgContainer.Terminate(suite.ctx)
		suite.NoError(err)
	}
}

func (suite *DatabaseTestSuite) TestPostgresClientRunDBMigration() {
	pool, err := suite.postgresClient.NewPostgresClient(suite.ctx)
	suite.NoError(err)
	suite.NotNil(pool)

	err = suite.postgresClient.PingDB(suite.ctx)
	suite.NoError(err)

	err = suite.postgresClient.RunDBMigration(suite.migrationURL)
	suite.NoError(err)

	err = suite.postgresClient.RunDBMigration("invalid-url://invalid")
	suite.Error(err)
}

func (suite *DatabaseTestSuite) TestSetupDatabaseSuccess() {
	dbPool, err := SetupDatabase(suite.ctx, suite.config, suite.migrationURL)
	suite.NoError(err)
	suite.NotNil(dbPool)
}

func (suite *DatabaseTestSuite) TestSetupDatabaseInvalidDSN() {
	invalidConfig := suite.config
	invalidConfig.DSN = "invalid_dsn"

	_, err := SetupDatabase(suite.ctx, invalidConfig, suite.migrationURL)
	suite.Error(err)
}

func (suite *DatabaseTestSuite) TestSetupDatabaseWithCanceledContext() {
	ctx, cancel := context.WithCancel(suite.ctx)
	cancel()

	_, err := SetupDatabase(ctx, suite.config, suite.migrationURL)
	suite.Error(err)
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
