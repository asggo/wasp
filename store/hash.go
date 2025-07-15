package store

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/text/unicode/norm"
)

const (
	saltSize = 16
	keySize  = 32
)

// ----------------------------------------------------------------------------
// passwordHash Struct
// ----------------------------------------------------------------------------
// passwordHash holds the information needed to construct a User's password
// Hash.
type passwordHash struct {
	Version uint8
	Memory  uint32
	Time    uint32
	Threads uint8
	Salt    [saltSize]byte
	Key     [keySize]byte
}

// verify takes a password and verifies it produces the same key.
func (p *passwordHash) verify(pwd string) bool {
	// Normalize our passphrase
	pwd = norm.NFKD.String(pwd)

	// Derive our hash.
	derived := argon2.IDKey([]byte(pwd), p.Salt[:], p.Time, p.Memory, p.Threads, keySize)

	// Compare the derived key to the existing key
	cmp := subtle.ConstantTimeCompare([]byte(p.Key), []byte(derived))

	return cmp == 1
}

// newArgonHash creates a new argonHash with the given parameters and a
// random salt.
func newPasswordHash(m, t uint32, p uint8, pwd string) (passwordHash, error) {
	var salt [saltSize]byte
	var hash passwordHash

	// Get a random salt value.
	_, err := rand.Read(salt[:])
	if err != nil {
		return hash, fmt.Errorf("could not NewPasswordHash: %v", err)
	}

	// Calculate the Key
	// Normalize our password
	pwd = norm.NFKD.String(pwd)

	// Derive our hash.
	key := argon2.IDKey([]byte(pwd), salt[:], t, m, p, keySize)

	hash = passwordHash{
		Version: 19,
		Time:    t,
		Memory:  m,
		Threads: p,
		Salt:    salt,
		Key:     key
	}

	return hash, nil
}

// newPasswordHashFromBytes creates a new passwordHash object from a JSON byte
// array.
func newPasswordHashFromBytes(data []byte) (passwordHash, error) {
	var ph passwordHash

	err := json.Unmarshal(data, &ph)
	if err != nil {
		return ph, fmt.Errorf("could not newPasswordHashFromBytes: %v", err)
	}

	return user, nil
}

//----------------------------------------------------------------------------
// passwordHash Storage Methods
//----------------------------------------------------------------------------
// getUserPasswordHash returns the passwordHash associated with the given
// userToken.
func (s *Store) getUserPasswordHash(ut userToken) (passwordHash, error) {
	var p passwordHash

	key := fmt.Sprintf(authHashKey, ut.String())
	data := s.read(authBucket, key)
	if data == nil {
		return p, fmt.Errorf("could not Store.GetUserPasswordHash: %v", err)
	}

	t, err := newPasswordHashFromBytes(data)
	if err != nil {
		return p, fmt.Errorf("could not Store.GetUserPasswordHash: %v", err)
	}

	return p, nil
}

// deleteUserPasswordHash deletes the passwordHash associated with the given
// userToken.
func (s *Store) deleteUserPasswordHash(ut userToken) error {
	key := fmt.Sprintf(authHashKey, ut.String())

	return s.Delete(authBucket, key)
}

// VerifyHash verifies a given passphrase generates the given hash.
func (s *Store) verifyHash(ut userToken, passphrase string) bool {
	ph, err := s.getUserPasswordHash(ut)
	if err != nil {
		return false
	}

	return ph.verify(passphrase)
}
