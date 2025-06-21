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

// argonDerive derives an Argon2id hash based on the stored parameters.
type argonHash struct {
	memory  uint32
	time    uint32
	threads uint8
	salt    [saltSize]byte
}

// derive takes a passphrase and returns a salt and derived key.
func (a argonHash) derive(passphrase string) string {
	// Normalize our passphrase
	passphrase = norm.NFKD.String(passphrase)

	// Derive our hash.
	key := argon2.IDKey([]byte(passphrase), a.salt[:], a.time, a.memory, a.threads, keySize)

	// Encode our data
	salt := base64.RawStdEncoding.EncodeToString(a.salt[:])
	hash := base64.RawStdEncoding.EncodeToString(key)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		a.memory, a.time, a.threads,
		salt, hash,
	)
}

// newArgonHash creates a new argonHash with the given parameters and salt.
func newArgonHash(m, t uint32, p uint8, salt [saltSize]byte) argonHash {
	return argonHash{
		time:    t,
		memory:  m,
		threads: p,
		salt:    salt,
	}
}

// newArgonHashFromString returns an argonHash from the salt and parameters
// extracted from the given string.
func newArgonHashFromString(hash string) (argonHash, error) {
	var ah argonHash
	var saltBytes []byte

	if !strings.HasPrefix(hash, "$argon2id$v=19") {
		return ah, fmt.Errorf("could not newArgonHashFromString: invalid hash type")
	}

	// Split the hash into its parts.
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return ah, fmt.Errorf("could not newArgonHashFromString: invalid hash split")
	}

	// Extract parameters
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &ah.memory, &ah.time, &ah.threads)
	if err != nil {
		return ah, fmt.Errorf("could not newArgonHashFromString: invalid parameters")
	}

	// Extract salt
	saltBytes, _ = base64.RawStdEncoding.DecodeString(parts[4])
	if len(saltBytes) != saltSize {
		return ah, fmt.Errorf("could not newArgonHashFromString: invalid salt length")
	}

	copy(ah.salt[:], saltBytes[:saltSize])

	return ah, nil
}

// GenerateHash creates a new Argon2id hash with the given passphrase.
func GenerateHash(passphrase string) (string, error) {
	var saltBytes [saltSize]byte
	var hash string

	// Get a random salt value.
	_, err := rand.Read(saltBytes[:])
	if err != nil {
		return hash, fmt.Errorf("could not argonHash.derive: %v", err)
	}

	// Create an argonHash with the SECOND RECOMMENDED option from RFC9106 and
	// derive our hash.
	argon := newArgonHash(64*1024, 4, 3, saltBytes)
	hash = argon.derive(passphrase)

	return hash, nil
}

// VerifyHash verifies a given passphrase generates the given hash.
func VerifyHash(hash, passphrase string) bool {
	argon, err := newArgonHashFromString(hash)
	if err != nil {
		return false
	}

	derived := argon.derive(passphrase)
	cmp := subtle.ConstantTimeCompare([]byte(hash), []byte(derived))

	return cmp == 1
}
