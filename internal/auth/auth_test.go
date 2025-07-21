package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// test hash password
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

// test make jwt
func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "test-secret", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	if token == "" {
		t.Fatalf("Token is empty")
	}
}

// test validate jwt
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "test-secret", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, "test-secret")
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("Invalid user ID: %v", validatedUserID)
	}
}

// test expired token
func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "test-secret", -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, "test-secret")
	if err == nil {
		t.Fatalf("Expected error for expired token")
	}
}

// test incorrect secret
func TestIncorrectSecret(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "test-secret", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, "incorrect-secret")
	if err == nil {
		t.Fatalf("Expected error for incorrect secret")
	}
}

// test invalid token
func TestInvalidToken(t *testing.T) {
	_, err := ValidateJWT("invalid-token", "test-secret")
	if err == nil {
		t.Fatalf("Expected error for invalid token")
	}
}