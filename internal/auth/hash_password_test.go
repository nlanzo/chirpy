package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if hashedPassword == password {
		t.Errorf("Hashed password is the same as the original password")
	}

	if err := CheckPasswordHash(password, hashedPassword); err != nil {
		t.Errorf("Error checking password hash: %v", err)
	}
}
