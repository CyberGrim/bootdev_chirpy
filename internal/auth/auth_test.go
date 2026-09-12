package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var correctPassword string = "CorrectPassword123"
var wrongPassword string = "InorrectPassword123"

func TestHashPassword(t *testing.T) {
	response, err := HashPassword(correctPassword)
	if err != nil {
		t.Fatal("Hash Creation Failed")
	}
	if response == correctPassword {
		t.Fatal("Response from Hash Creation matched test password")
	}
	if response == "" {
		t.Fatal("Responce from Hash Creation was empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPassword(correctPassword)

	if err != nil {
		t.Fatal("Hash Creation Failed")
	}

	responseTruth, err := CheckPasswordHash(correctPassword, hash)
	if err != nil {
		t.Fatal("Hash Check Function Failed")
	}
	if responseTruth == false {
		t.Fatal("Hash Check returned false - Expected true")
	}

	responseFalse, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatal("Hash Check Function Failed")
	}
	if responseFalse == true {
		t.Fatal("Hash Check returned true - Expected false")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	gotUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("got user ID %v, want %v", gotUserID, userID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if _, err := ValidateJWT(token, "test-secret"); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if _, err := ValidateJWT(token, "wrong-secret"); err == nil {
		t.Fatal("expected token signed with the wrong secret to be rejected")
	}
}

func TestValidateJWTRejectsMalformedToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{name: "empty string", token: ""},
		{name: "wrong segment count", token: "not.a.token"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ValidateJWT(test.token, "test-secret"); err == nil {
				t.Fatalf("expected malformed token %q to be rejected", test.token)
			}
		})
	}
}

func TestCheckPasswordHashRejectsMalformedHash(t *testing.T) {
	ok, err := CheckPasswordHash(correctPassword, "")
	if err == nil {
		t.Fatal("expected malformed hash to return an error")
	}
	if ok {
		t.Fatal("expected malformed hash not to match")
	}
}
