package store

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

var (
	testDeriveMemory         = 64 * 1024
	testDeriveTime           = 4
	testDeriveThreads        = 3
	testDeriveGoodPassword   = "password"
	testDeriveBadAlgHash     = "$argon2i$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	testDeriveBadVersionHash = "$argon2id$v=18$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	testDeriveBadSplitHash   = "$argon2id$v=19$m=65536,t=4,p=3c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	testDeriveBadParamsHash  = "$argon2id$v=19$m=65536t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	testDeriveBadSaltHash    = "$argon2id$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2Fsd$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
)

func TestDeriver(t *testing.T) {
	t.Run("Test newPasswordHash", testNewPasswordHash)
	t.Run("Test newArgonHashFromString", testNewArgonHashFromString)
}

func testNewPasswordHash(t *testing.T) {
	fmt.Println(t.Name())

	// Create two argonHash objects.
	ah, err := newPasswordHash(testDeriveMemory, testDeriveTime, testDeriveThreads)
	ah2, err := newPasswordHash(testDeriveMemory, testDeriveTime, testDeriveThreads)

	// They should each have different salts.
	if bytes.Equal(ah.salt[:], ah1.salt[:]) {
		t.Fatalf("Expected unique salts, received %+v, %+v", ah.salt, ah1.salt)
	}

	// The derived hashes should be different as well
	hash1 := ah.derive(testDeriveGoodPassword)
	hash2 := ah2.derive(testDeriveGoodPassword)
	if hash1 == hash2 {
		t.Fatal("Expected unique hashes, received", hash1, hash2)
	}
}

func testNewArgonHashFromString(t *testing.T) {
	fmt.Println(t.Name())

	_, err := newArgonHashFromString(testDeriveBadAlgHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash type") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(testDeriveBadVersionHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash type") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(testDeriveBadSplitHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash split") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(testDeriveBadParamsHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid parameters") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(testDeriveBadSaltHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid salt length") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	a, err := newArgonHashFromString(testDeriveGoodHash)
	if err != nil {
		t.Fatal("Expected no error received", err)
	}

	if a.memory != testDeriveMemory {
		t.Fatal("Expected", testDeriveMemory, "received", a.memory)
	}

	if a.time != testDeriveTime {
		t.Fatal("Expected", testDeriveTime, "received", a.time)
	}

	if a.threads != testDeriveThreads {
		t.Fatal("Expected", testDeriveThreads, "received", a.threads)
	}
}
