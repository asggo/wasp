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
}

// derive takes a passphrase and returns a salt and derived key.
func (p *passwordHash) derive(pwd string) string {
	// Normalize our passphrase
	passphrase = norm.NFKD.String(passphrase)

	// Derive our hash.
	key := argon2.IDKey([]byte(passphrase), p.Salt[:], p.Time, p.Memory, p.Threads, keySize)

	// Encode our data
	salt := base64.RawStdEncoding.EncodeToString(p.Salt[:])
	hash := base64.RawStdEncoding.EncodeToString(key)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		p.Version, p.Memory, p.Time, p.Threads,
		salt, hash,
	)
}

// newArgonHash creates a new argonHash with the given parameters and a
// random salt.
func newPasswordHash(m, t uint32, p uint8) (passwordHash, error) {
	var salt [saltSize]byte
	var hash passwordHash

	// Get a random salt value.
	_, err := rand.Read(salt[:])
	if err != nil {
		return hash, fmt.Errorf("could not NewPasswordHash: %v", err)
	}

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
func (s *Store) verifyHash(ut userToken, hash, passphrase string) bool {
	key := fmt.Sprintf(userHashKey, ut.String())
	data := s.read(userBucket, key)
	if data == nil {
		return false
	}

	ah, err := newArgonHashFromString(string(data))
	if err != nil {
		return false
	}

	derived := ah.derive(passphrase)
	cmp := subtle.ConstantTimeCompare([]byte(hash), []byte(derived))

	return cmp == 1
}
