package store

import (
	"fmt"
	"strings"
	"testing"
)

const (
	tokenBase32Size = 52
)

func TestTokens(t *testing.T) {
	t.Run("Test UserToken", testUserToken)
	t.Run("Test SessionToken", testSessionToken)
}

func testTokenBytes(t *testing.T, s string) {
	data, err := tokenEncoder.DecodeString(s)
	if err != nil {
		t.Fatal("Could not parse token bytes:", err)
	}

	if len(data) != tokenSize {
		t.Fatal("Expected", tokenSize, "bytes, received", len(data))
	}
}

func testUserToken(t *testing.T) {
	fmt.Println(t.Name())

	token := NewUserToken().String()

	if !strings.HasPrefix(token, userTokenPrefix) {
		t.Fatal("UserToken has incorrect prefix.")
	}

	parsed, err := parseUserToken(token)
	if err != nil {
		t.Fatal("Expected no error, recieved", err)
	}

	if parsed.String() != token {
		t.Fatal("Expected", token, ", received", parsed.String())
	}

	token = strings.TrimPrefix(token, userTokenPrefix)
	if len(token) != tokenBase32Size {
		t.Fatal("Expected", tokenBase32Size, "base32 characters, received", len(token))
	}

	testTokenBytes(t, token)
}

func testSessionToken(t *testing.T) {
	fmt.Println(t.Name())

	token := NewSessionToken().String()

	if !strings.HasPrefix(token, sessionTokenPrefix) {
		t.Fatal("SessionToken has incorrect prefix.")
	}

	parsed, err := parseSessionToken(token)
	if err != nil {
		t.Fatal("Expected no error, recieved", err)
	}

	if parsed.String() != token {
		t.Fatal("Expected", token, ", received", parsed.String())
	}

	token = strings.TrimPrefix(token, sessionTokenPrefix)
	if len(token) != tokenBase32Size {
		t.Fatal("Expected", tokenBase32Size, "base32 characters, received", len(token))
	}

	testTokenBytes(t, token)
}
