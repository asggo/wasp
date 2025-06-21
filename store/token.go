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
)

// tokenEncoder is used to encoded and decode our tokens using a standard
// Base32 encoder with no padding.
var tokenEncoder = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)

//----------------------------------------------------------------------------
// UserToken
//----------------------------------------------------------------------------

// UserToken represents a user token.
type UserToken [tokenSize]byte

// String converts a UserToken object to a string.
func (u UserToken) String() string {
	token := tokenEncoder.EncodeToString(u[:])

	return fmt.Sprintf("%s%s", userTokenPrefix, token)
}

// NewUserToken generates a random UserToken.
func NewUserToken() UserToken {
	var ut UserToken

	bytes := newTokenBytes()
	copy(ut[:], bytes[:])

	return ut
}

// parseUserToken takes a string in the form of user_base32 and parses it
// into an UserToken
func parseUserToken(s string) (UserToken, error) {
	var ut UserToken

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
// SessionToken
//----------------------------------------------------------------------------

// SessionToken represents a session token.
type SessionToken [tokenSize]byte

// String converts a SessionToken object to a string.
func (s SessionToken) String() string {
	token := tokenEncoder.EncodeToString(s[:])

	return fmt.Sprintf("%s%s", sessionTokenPrefix, token)
}

// NewSessionToken generates a random SessionToken.
func NewSessionToken() SessionToken {
	var st SessionToken

	bytes := newTokenBytes()
	copy(st[:], bytes[:])

	return st
}

// parseSessionToken takes a string in the form of sess_base32 and parses it
// into an SessionToken
func parseSessionToken(s string) (SessionToken, error) {
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
