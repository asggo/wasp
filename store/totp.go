package store

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

//----------------------------------------------------------------------------
// Totp Storage Methods
//----------------------------------------------------------------------------

// GetTotpUrl returns the URL for the given user's Totp secret in Google Auth
// format.
func (s *Store) GetTotpUrl(ut UserToken) (string, error) {
	key := fmt.Sprintf(userTotpKey, ut.String())

	data := s.read(userBucket, key)
	if data == nil {
		return "", fmt.Errorf("could not Store.GetTotpUrl: %v", err)
	}

	secret := strings.TrimPrefix(totpTokenPrefix, string(data))
	otpStr := "otpauth://totp/%s?secret=%s&issuer=%s"
	url := fmt.Sprintf(otpStr, s.cfg.TotpName, secret, s.cfg.TotpIssuer)

	return url, nil
}

// VerifyTotp returns true if the given totp value matches the value generated
// for the given user.
func (s *Store) verifyTotp(ut UserToken, expected string) bool {
	key := fmt.Sprintf(userTotpKey, ut.String())

	data := s.read(userBucket, key)
	if data == nil {
		return false
	}

	tt, err := parseTotpToken(string(data))
	if err != nil {
		return false
	}

	ts := time.Now().Unix()
	derived := generateSha256Totp(
		tt[:],
		ts,
		s.cfg.TotpLength,
		s.cfg.TotpStart,
		s.cfgTotpStep,
	)

	cmp := subtle.ConstantTimeCompare([]byte(derived), []byte(expected))

	return cmp == 1
}

//----------------------------------------------------------------------------
// Totp Helper Functions
//----------------------------------------------------------------------------

// generateSha256Totp creates a Time-based OTP value as defined in RFC 6238.
//
// This method is a Golang port of the Java implentation defined in
// https://datatracker.ietf.org/doc/html/rfc6238#appendix-A.
func generateSha256Totp(secret []byte, start, step, ts int64, l int) string {
	// Calculate the number of steps since totpStart, and convert it to an
	// 8-byte array.
	steps := uint64(math.Floor(float64((ts - start) / step)))
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, steps)

	// Calculate the HMAC using the given key.
	mac := hmac.New(sha256.New, secret)
	mac.Write(buf)
	h := mac.Sum(nil)

	// Do dynamic truncation
	// The offset is the last four bits of the last byte in the hash.
	off := h[len(h)-1] & 0xf

	// The binary value is the four bytes starting at the offset with the
	// highest bit of the first byte set to 0.
	bin := h[off : off+4]
	bin[0] = bin[0] & 0x7f

	// Get the OTP value
	// Convert the binary value to an integer and calculate the modulus based
	// on the number of digits desired in the OTP.
	i := binary.BigEndian.Uint32(bin)

	switch l {
	case 8:
		return fmt.Sprintf("%08d", i%100000000)
	case 7:
		return fmt.Sprintf("%07d", i%10000000)
	default:
		return fmt.Sprintf("%06d", i&1000000)
	}
}
