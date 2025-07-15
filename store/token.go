package store

// Tokens are always randomly generated and only used as identifiers. All
// tokens have the form prefix_base32encodedbytes.
import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

const (
	tokenSize          = 32
	userTokenPrefix    = "user_"
	sessionTokenPrefix = "sess_"
	totpTokenPrefix    = "totp_"
)

// tokenEncoder is used to encoded and decode our tokens using a standard
// Base32 encoder with no padding.
var tokenEncoder = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)

//----------------------------------------------------------------------------
// userToken
//----------------------------------------------------------------------------

// userToken represents a user token.
type userToken [tokenSize]byte

// String converts a UserToken object to a string.
func (u userToken) String() string {
	token := tokenEncoder.EncodeToString(u[:])

	return fmt.Sprintf("%s%s", userTokenPrefix, token)
}

// newUserToken generates a random userToken.
func newUserToken() userToken {
	var ut UserToken

	bytes := newTokenBytes()
	copy(ut[:], bytes[:])

	return ut
}

// parseUserToken takes a string in the form of user_base32 and parses it
// into an userToken
func parseUserToken(s string) (userToken, error) {
	var ut userToken

	if !strings.HasPrefix(s, userTokenPrefix) {
		return ut, fmt.Errorf("could not parseUserToken: invalid prefix")
	}

	s = strings.TrimPrefix(s, userTokenPrefix)

	data, err := tokenEncoder.DecodeString(s)
	if err != nil {
		return ut, fmt.Errorf("could not parseUserToken: %v", err)
	}

	if len(data) != tokenSize {
		return ut, fmt.Errorf("could not parseUserToken: invalid length")
	}

	copy(ut[:], data)

	return ut, nil
}

//----------------------------------------------------------------------------
// sessionToken
//----------------------------------------------------------------------------

// sessionToken represents a session token.
type sessionToken [tokenSize]byte

// String converts a sessionToken object to a string.
func (s sessionToken) String() string {
	token := tokenEncoder.EncodeToString(s[:])

	return fmt.Sprintf("%s%s", sessionTokenPrefix, token)
}

// NewSessionToken generates a random sessionToken.
func newSessionToken() sessionToken {
	var st sessionToken

	bytes := newTokenBytes()
	copy(st[:], bytes[:])

	return st
}

// parseSessionToken takes a string in the form of sess_base32 and parses it
// into an sessionToken
func parseSessionToken(s string) (sessionToken, error) {
	var st SessionToken

	if !strings.HasPrefix(s, sessionTokenPrefix) {
		return st, fmt.Errorf("could not parseSessionToken: invalid prefix")
	}

	s = strings.TrimPrefix(s, sessionTokenPrefix)

	data, err := tokenEncoder.DecodeString(s)
	if err != nil {
		return st, fmt.Errorf("could not parseSessionToken: %v", err)
	}

	if len(data) != tokenSize {
		return st, fmt.Errorf("could not parseSessionToken: invalid length")
	}

	copy(st[:], data)

	return st, nil
}

//----------------------------------------------------------------------------
// totpToken
//----------------------------------------------------------------------------

// totpToken represents a TOTP token.
type totpToken [tokenSize]byte

// String converts a TotpToken object to a string.
func (t totpToken) String() string {
	token := tokenEncoder.EncodeToString(s[:])

	return fmt.Sprintf("%s%s", totpTokenPrefix, token)
}

// NewTotpToken generates a random totpToken.
func NewTotpToken() totpToken {
	var tt totpToken

	bytes := newTokenBytes()
	copy(tt[:], bytes[:])

	return tt
}

// parseTotpToken takes a string in the form of totp_base32 and parses it
// into an totpToken
func parseTotpToken(s string) (totpToken, error) {
	var tt totpToken

	if !strings.HasPrefix(s, totpTokenPrefix) {
		return tt, fmt.Errorf("could not parseTotpToken: invalid prefix")
	}

	s = strings.TrimPrefix(s, totpTokenPrefix)

	data, err := tokenEncoder.DecodeString(s)
	if err != nil {
		return tt, fmt.Errorf("could not parseTotpToken: %v", err)
	}

	if len(data) != tokenSize {
		return tt, fmt.Errorf("could not parseTotpToken: invalid length")
	}

	copy(tt[:], data)

	return tt, nil
}

//----------------------------------------------------------------------------
// Helper functions
//----------------------------------------------------------------------------

// newTokenBytes returns a byte slice with tokenSize random bytes.
func newTokenBytes() [tokenSize]byte {
	var bytes [tokenSize]byte

	_, err := rand.Read(bytes[:])
	if err != nil {
		panic(fmt.Errorf("Could not newTokenBytes: %v", err))
	}

	return bytes
}
