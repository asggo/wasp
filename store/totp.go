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

// ----------------------------------------------------------------------------
// totp Stuct
// ----------------------------------------------------------------------------
// totp holds a User's Totp object.
type totp struct {
	Secret TotpToken `json:"secret"`
	Length int64     `json:"length"`
	Start  int64     `json:"start"`
	Step   int64     `json:"step"`
	Name   string    `json:"name"`
	Issuer string    `json:"issuer"`
}

// bytes converts a Totp object to a JSON byte array.
func (t *totp) bytes() ([]byte, error) {
	var b []byte

	b, err := json.Marshal(s)
	if err != nil {
		return b, fmt.Errorf("could not Totp.bytes: %v", err)
	}

	return b, nil
}

// GetCode creates a Time-based OTP value as defined in RFC 6238.
//
// This method is a Golang port of the Java implentation defined in
// https://datatracker.ietf.org/doc/html/rfc6238#appendix-A.
func (t *totp) getCode(ts int64) string {
	// Calculate the number of steps since totpStart, and convert it to an
	// 8-byte array.
	steps := uint64(math.Floor(float64((ts - t.Start) / t.Step)))
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

	switch t.Length {
	case 8:
		return fmt.Sprintf("%08d", i%100000000)
	case 7:
		return fmt.Sprintf("%07d", i%10000000)
	default:
		return fmt.Sprintf("%06d", i&1000000)
	}
}

func (t *totp) url() string {
	secret := strings.TrimPrefix(totpTokenPrefix, t.Secret.String())
	otpStr := "otpauth://totp/%s?secret=%s&issuer=%s"
	url := fmt.Sprintf(otpStr, t.Name, secret, t.Issuer)
}

// NewTotp returns a new Totp object.
func newTotp(length, start, step int64, name, issuer string) totp {
	return Totp{
		Secret: NewTotpToken(),
		Length: length,
		Start:  start,
		Step:   step,
		Name:   name,
		Issuer: issuer,
	}
}

// newTotpFromBytes creates a new Totp object from a JSON byte array.
func newTotpFromBytes(data []byte) (totp, error) {
	var t Totp

	err := json.Unmarshal(data, &t)
	if err != nil {
		return t, fmt.Errorf("could not NewTotpFromBytes: %v", err)
	}

	return t, nil
}

//----------------------------------------------------------------------------
// Totp Storage Methods
//----------------------------------------------------------------------------

// GetTotpUrl returns the URL for the given user's Totp secret in Google Auth
// format.
func (s *Store) GetUserTotpUrl(ut userToken) (string, error) {
	t, err := s.getUserTotp(ut)
	if err != nil {
		return "", fmt.Errorf("could not Store.GetUserTotpUrl: %v", err)
	}

	return t.url(), nil
}

// getUserTotp returns the TotpToken associated with the given UserToken.
func (s *Store) getUserTotp(ut userToken) (totp, error) {
	var t Totp

	key := fmt.Sprintf(authTotpKey, ut.String())
	data := s.read(authBucket, key)
	if data == nil {
		return t, fmt.Errorf("could not Store.GetUserTotp: %v", err)
	}

	t, err := newTotpFromBytes(data)
	if err != nil {
		return t, fmt.Errorf("could not Store.GetUserTotp: %v", err)
	}

	return t, nil
}

// getRecentTotpCode retrieves the last successfully used totp code from the
// store.
func (s *Store) getRecentTotpCode(ut userToken) string {
	key := fmt.Sprintf(authRecentCodeKey, ut.String())

	return string(s.read(authBucket, key))
}

// saveRecentTotpCode saves the most recent successfully used totp code to the
// store.
func (s *Store) saveRecentTotpCode(ut userToken, code string) {
	key := fmt.Sprintf(authUsedCodeKey, ut.String())
	s.write(authBucket, key, []byte(code))
}

// deleteUserTotpToken deletes the TotpToken associated with the given
// UserToken.
func (s *Store) deleteUserTotp(ut userToken) error {
	key := fmt.Sprintf(authTotpKey, ut.String())

	return s.Delete(authBucket, key)
}

// verifyTotp returns true if the expected totp value matches the value
// generated for the given user.
func (s *Store) verifyTotp(ut userToken, expected string) bool {
	totp, err := s.getUserTotp(ut)
	if err != nil {
		return false
	}

	// Get the current code and compare it to the expected code.
	now := time.Now().Unix()
	curr := totp.getCode(now)

	cmp := subtle.ConstantTimeCompare([]byte(curr), []byte(expected))
	if cmp == 0 {
		// If the comparison fails, calculate the previous code and test it.
		prev := totp.getCode(now - s.cfg.Step)

		return subtle.ConstantTimeCompare([]byte(prev), []byte(expected)) == 1
	} else {
		return true
	}
}
