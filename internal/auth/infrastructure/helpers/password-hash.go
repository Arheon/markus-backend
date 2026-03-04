package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	// CompareHashAndPassword compares the input password with the stored hash
	// it extracts the salt and performs the hash internally
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
