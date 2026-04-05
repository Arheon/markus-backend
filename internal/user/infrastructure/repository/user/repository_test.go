package user

import (
	"context"
	"testing"
	"time"

	sharedHelpers "github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	domainEntity "github.com/Arheon/markus-backend/internal/user/domain/entity"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	testContainersPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	container *testContainersPostgres.PostgresContainer
	db        *gorm.DB
	repo      repository.UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	pgContainer, err := testContainersPostgres.Run(ctx,
		"postgres:15-alpine",
		testContainersPostgres.WithDatabase("testdb"),
		testContainersPostgres.WithUsername("user"),
		testContainersPostgres.WithPassword("pass"),
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

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.container.Terminate(context.Background())
}

func (s *UserRepositoryTestSuite) SetupTest() {
	s.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func (s *UserRepositoryTestSuite) TestCreateUserOrErrorIfExists() {
	ctx := context.Background()

	userName := "arheon"
	password := "password"
	passwordHash, err := sharedHelpers.HashPassword(password)
	s.NoError(err)

	tests := []struct {
		name      string
		username  string
		password  string
		prepareDB func()
		wantName  string
		err       error
	}{
		{
			name:      "Success create user",
			username:  userName,
			password:  passwordHash,
			prepareDB: nil,
			wantName:  userName,
			err:       nil,
		},

		{
			name:     "User already exists",
			username: userName,
			password: passwordHash,
			prepareDB: func() {
				s.repo.CreateNewUserOrErrorIfExists(ctx, userName, passwordHash)
			},
			wantName: "",
			err:      ErrUserIsExists,
		},
	}

	s.NotNil(s.repo, "Репозиторий не должен быть nil")
	s.NotNil(s.db, "Объект БД не должен быть nil")
	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			if tt.prepareDB != nil {
				tt.prepareDB()
			}

			user, err := s.repo.CreateNewUserOrErrorIfExists(ctx, tt.username, tt.password)

			if tt.err != nil {
				// Мы ждем ошибку
				s.ErrorIs(err, tt.err, "Должна быть ошибка: %v", tt.err)
				return // Дальше проверять user.Username нельзя, user будет nil
			}

			// Мы НЕ ждем ошибку
			s.NoError(err)
			s.NotNil(user)
			s.Equal(tt.wantName, user.Username)
		})
	}
}

func TestRepository_CreateNewUserOrErrorIfExists(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
