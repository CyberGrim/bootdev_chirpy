package auth

import (
	"testing"
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
