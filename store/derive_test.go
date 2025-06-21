package store

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

var (
	goodPassword   = "password"
	goodSalt       = "c2FsdHNhbHRzYWx0c2FsdA"
	goodHash       = "$argon2id$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badAlgHash     = "$argon2i$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badVersionHash = "$argon2id$v=18$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badSplitHash   = "$argon2id$v=19$m=65536,t=4,p=3c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badParamsHash  = "$argon2id$v=19$m=65536t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badSaltHash    = "$argon2id$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2Fsd$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Zds"
	badKeyHash     = "$argon2id$v=19$m=65536,t=4,p=3$c2FsdHNhbHRzYWx0c2FsdA$gqlSqdyZN20rXTgEHcac5Bfuege+r2YT0FKkcRK7Ztt"
)

func TestDeriver(t *testing.T) {
	t.Run("Test newArgonHash", testNewArgonHash)
	t.Run("Test newArgonHashFromString", testNewArgonHashFromString)
	t.Run("Test GenerateHash", testGenerateHash)
	t.Run("Test VerifyHash", testVerifyHash)
}

func testNewArgonHash(t *testing.T) {
	fmt.Println(t.Name())

	var saltBytes [saltSize]byte

	bytes, _ := base64.RawStdEncoding.DecodeString(goodSalt)
	copy(saltBytes[:], bytes)

	argon := newArgonHash(64*1024, 4, 3, saltBytes)
	hash := argon.derive(goodPassword)

	if hash != goodHash {
		t.Fatal("Expected", goodHash, ", received", hash)
	}

	argon = newArgonHash(64*1024, 4, 2, saltBytes)
	hash = argon.derive(goodPassword)

	if hash == goodHash {
		t.Fatal("Did not expect", goodHash, "to equal", hash)
	}
}

func testNewArgonHashFromString(t *testing.T) {
	fmt.Println(t.Name())

	_, err := newArgonHashFromString(badAlgHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash type") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(badVersionHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash type") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(badSplitHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid hash split") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(badParamsHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid parameters") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	_, err = newArgonHashFromString(badSaltHash)
	if (err == nil) || !strings.Contains(err.Error(), "invalid salt length") {
		t.Fatal("Expected invalid hash type, received", err)
	}

	a, err := newArgonHashFromString(goodHash)
	if err != nil {
		t.Fatal("Expected no error received", err)
	}

	if (a.memory != 65536) || (a.time != 4) || (a.threads != 3) {
		t.Fatal("Expected 65536, 4, 3, received", a.memory, a.time, a.threads)
	}

	salt := base64.RawStdEncoding.EncodeToString(a.salt[:])
	if salt != goodSalt {
		t.Fatal("Expected", goodSalt, ", received", salt)
	}
}

func testGenerateHash(t *testing.T) {
	fmt.Println(t.Name())

	hash1, err := GenerateHash(goodPassword)
	if err != nil {
		t.Fatal("Expected no error, recieved", err)
	}

	hash2, _ := GenerateHash(goodPassword)
	if hash1 == hash2 {
		t.Fatal("Did not generate unique hashes")
	}
}

func testVerifyHash(t *testing.T) {
	fmt.Println(t.Name())

	if !VerifyHash(goodHash, goodPassword) {
		t.Fatal("Expected hashes to match, but they did not")
	}

	if VerifyHash(badKeyHash, goodPassword) {
		t.Fatal("Expected hash to not match, but they did")
	}
}
