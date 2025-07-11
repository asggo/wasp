package store

// The seed and test vectors are defined in the Java implementation in the
// TOTP RFC, https://datatracker.ietf.org/doc/html/rfc6238#appendix-A

import (
	"fmt"
	"testing"
)

var (
	testTotpSha256Seed = []byte{
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
		0x31, 0x32,
	}
)

type testTotpValue struct {
	time int64
	totp string
}

func TestGenerateSha256Totp(t *testing.T) {
	fmt.Println(t.Name())

	tests := []testTotpValue{
		{59, "46119246"},
		{1111111109, "68084774"},
		{1111111111, "67062674"},
		{1234567890, "91819424"},
		{2000000000, "90698825"},
		{20000000000, "77737706"},
	}

	for _, test := range tests {
		have := GenerateSha256Totp(testTotpSha256Seed, test.time, 8)
		if have != test.totp {
			t.Fatal("Expected", test.totp, ", received", have)
		}
	}
}
