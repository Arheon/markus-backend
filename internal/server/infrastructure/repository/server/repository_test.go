package server

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	domainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	serverDomainRepository "github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	testContainerPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServerRepositoryTestSuite struct {
	suite.Suite
	container *testContainerPostgres.PostgresContainer
	db        *gorm.DB
	repo      serverDomainRepository.ServerRepository
}

func (s *ServerRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	pgContainer, err := testContainerPostgres.Run(ctx,
		"postgres:15-alpine",
		testContainerPostgres.WithDatabase("testdb"),
		testContainerPostgres.WithUsername("user"),
		testContainerPostgres.WithPassword("pass"),
		// Это заставит Testcontainers ждать, пока база реально ответит на запрос,
		// а не просто пока запустится процесс в Docker
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2). // Важно: Postgres пишет это дважды при старте
				WithStartupTimeout(30*time.Second),
		),
	)
	s.NoError(err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	s.NoError(err)

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{})
	s.NoError(err)

	err = db.AutoMigrate(domainEntity.User{})
	s.NoError(err)

	s.container = pgContainer
	s.db = db
	s.repo = NewRepository(db)
}

func (s *ServerRepositoryTestSuite) TearDownSuite() {
	s.container.Terminate(context.Background())
}

func (s *ServerRepositoryTestSuite) SetupTest() {
	s.db.Exec("TRUNCATE TABLE servers RESTART IDENTITY CASCADE")
	s.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func (s *ServerRepositoryTestSuite) TestCreateServer() {
	ctx := context.Background()
	userId := "b4cf8687-a33c-43cc-a573-d2b8a931de08"
	userName := "arheon"
	password, err := helpers.HashPassword("password")
	serverName := "testServer"

	s.NoError(err)

	tests := []struct {
		name        string
		prepareDB   func() error
		expextedErr error
		serverName  *string
	}{
		{
			name:        "Error undefined user",
			prepareDB:   nil,
			expextedErr: ErrUndefinedUser,
			serverName:  &serverName,
		},
		{
			name: "Success create server",
			prepareDB: func() error {
				err := gorm.G[any](s.db).Exec(ctx, fmt.Sprintf("INSERT INTO users (id, username, password) VALUES ('%s', '%s', '%s')", userId, userName, password))
				return err
			},
			expextedErr: nil,
			serverName:  &serverName,
		},
		{
			name: "Error create server",
			prepareDB: func() error {
				err := gorm.G[any](s.db).Exec(ctx, fmt.Sprintf("INSERT INTO users (id, username, password) VALUES ('%s', '%s', '%s')", userId, userName, password))
				return err
			},
			expextedErr: errors.New("Undefined db"),
			serverName:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			s.SetupTest()
			if tt.prepareDB != nil {
				err := tt.prepareDB()
				s.NoError(err)
				return
			}

			server, err := s.repo.CreateNewServer(ctx, userId, *tt.serverName)
			if tt.expextedErr != nil {
				s.ErrorIs(err, tt.expextedErr)
				return
			}

			s.NoError(err)
			s.NotNil(server)
		})
	}
}

func TestRepository_Main(t *testing.T) {
	suite.Run(t, new(ServerRepositoryTestSuite))
}
