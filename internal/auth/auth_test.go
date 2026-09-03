package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	testID := uuid.New()
	token, err := MakeJWT(testID, "test1", time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT() error: %v", err)
	}
	compareID, err := ValidateJWT(token, "test1")
	if err != nil {
		t.Fatalf("ValidateJWT() error: %v", err)
	}
	if testID != compareID {
		t.Errorf("MakeJWT() = %q, expected %q", testID, compareID)
	}
}

func TestExpiredJWT(t *testing.T) {
	testID := uuid.New()
	token, err := MakeJWT(testID, "test1", -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT() error: %v", err)
	}
	_, err = ValidateJWT(token, "test1")
	if err == nil {
		t.Fatal("ValidateJWT() error: expired token to be rejected")
	}
}

func TestWrongSecret(t *testing.T) {
	testID := uuid.New()
	token, err := MakeJWT(testID, "test1", time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT() error: %v", err)
	}
	_, err = ValidateJWT(token, "test2")
	if err == nil {
		t.Fatal("ValidateJWT() error: expected validation to fail with wrong secret")
	}
}
