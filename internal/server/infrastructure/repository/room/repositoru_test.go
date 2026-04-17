package room

import (
	"context"
	"time"

	"github.com/Arheon/markus-backend/internal/server/domain/entity"
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	testContainerPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RoomRepositoryTestSuite struct {
	suite.Suite
	container *testContainerPostgres.PostgresContainer
	db        *gorm.DB
	repo      repository.RoomRepository
}

func (r *RoomRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	pgContainer, err := testContainerPostgres.Run(
		ctx,
		"postgres:15-alpine",
		testContainerPostgres.WithDatabase("testdb"),
		testContainerPostgres.WithUsername("user"),
		testContainerPostgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	r.NoError(err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	r.NoError(err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	r.NoError(err)

	r.container = pgContainer
	r.db = db
	r.repo = NewRepository(db)
}

func (r *RoomRepositoryTestSuite) TearDownSuite() {
	r.container.Terminate(context.Background())
}

func (r *RoomRepositoryTestSuite) SetupTest() {
	r.db.Exec("TRUNCATE TABLE rooms RESTART IDENTITY CASCADE")
	r.db.Exec("TRUNCATE TABLE room_category RESTART IDENTITY CASCADE")
}

func (r *RoomRepositoryTestSuite) TestCreateRoom() {
	ctx := context.Background()
	serverID := "ce11f5fe-11b8-4282-8a12-75d00c010a1a"
	serverName := "test"

	roomName := "test_room"

	server := &entity.Server{
		ID:   serverID,
		Name: serverName,
	}

	err := gorm.G[entity.Server](r.db).Create(ctx, server)
	r.NoError(err)

	tests := []struct {
		name        string
		prepareDB   func() error
		expextedErr error
	}{
		{
			name:        "Success create room",
			prepareDB:   nil,
			expextedErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		r.Run(tt.name, func() {
			r.SetupTest()
			if tt.prepareDB != nil {
				err := tt.prepareDB()
				r.NoError(err)
			}

			room, err := r.repo.CreateRoom(ctx, roomName, serverID, nil)
			if tt.expextedErr != nil {
				r.ErrorIs(err, tt.expextedErr)
			} else {
				r.NoError(err)
			}

			r.NotNil(room)
		})
	}
}
