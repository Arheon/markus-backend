package createuser

import (
	"context"
	"errors"
	"testing"

	"github.com/Arheon/markus-backend/internal/shared/domain/entity"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	domainEntity "github.com/Arheon/markus-backend/internal/user/domain/entity"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	repositoryMock "github.com/Arheon/markus-backend/internal/user/domain/repository/mocks"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCommand_Handle(t *testing.T) {
	userId := "4f953e99-12a7-474b-9b8b-506a3b14d34e"
	userName := "arheon"
	password := "password"
	passwordHash, _ := helpers.HashPassword(password)

	user := &domainEntity.User{
		User: entity.User{
			ID:       userId,
			Username: userName,
			Password: passwordHash,
		},
		Name: "",
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo func() repository.UserRepository
		// Named input parameters for target function.
		username string
		password string
		want     string
		wantErr  bool
	}{
		{
			name: "success user create",
			userRepo: func() repository.UserRepository {
				repo := repositoryMock.NewMockUserRepository(t)
				repo.On("CreateNewUserOrErrorIfExists",
					mock.Anything,
					"arheon",
					mock.MatchedBy(func(hash string) bool {
						err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

						return err == nil
					}),
				).Return(user, nil)

				return repo
			},
			username: userName,
			password: password,
			wantErr:  false,
			want:     userId,
		},
		{
			name: "user is exists",
			userRepo: func() repository.UserRepository {
				repo := repositoryMock.NewMockUserRepository(t)
				repo.On("CreateNewUserOrErrorIfExists",
					mock.Anything,
					"arheon",
					mock.MatchedBy(func(hash string) bool {
						err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

						return err == nil
					}),
				).Return(nil, errors.New("user is exists"))

				return repo
			},
			username: userName,
			password: password,
			wantErr:  true,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommand(context.Background(), tt.userRepo())
			got, gotErr := c.Handle(tt.username, tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Handle() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Handle() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
