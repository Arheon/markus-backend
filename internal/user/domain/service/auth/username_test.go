package auth

import (
	"errors"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		err      error
	}{
		{
			name:     "success username validation",
			username: "arheon",
			wantErr:  false,
			err:      nil,
		},
		{
			name:     "username less than 2 characters",
			username: "a",
			wantErr:  true,
			err:      ErrUsernameTooShort,
		},
		{
			name:     "username more than 32 characters",
			username: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantErr:  true,
			err:      ErrUsernameTooLong,
		},
		{
			name:     "username has control characters",
			username: "44\t44wwww",
			wantErr:  true,
			err:      ErrUsernameHasControlCharacters,
		},
		{
			name:     "username has invisible characters",
			username: "sdfsdfsdf\u200Dsdfsfsd",
			wantErr:  true,
			err:      ErrUsernameHasInvisibleCharacters,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateUsername(tt.username)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateUsername() failed: %v", gotErr)
				}

				if !errors.Is(tt.err, gotErr) {
					t.Errorf("ValidateUsername() failed with error: %v, but expects: %v", gotErr, tt.err)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("ValidateUsername() succeeded unexpectedly")
			}
		})
	}
}
