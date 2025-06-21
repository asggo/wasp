package store

import (
	"fmt"
)

const (
	maxFailCount = 10
)

var (
	failedKey = "%s:failed"
	hashKey   = "%s:hash"
)

// ----------------------------------------------------------------------------
// Authentication Storage Methods
// ----------------------------------------------------------------------------
// Authenticate takes a passphrase and verifies it matches the user's original
// passphrase.
func (s *Store) AuthenticateUser(ut UserToken, passphrase string) bool {
	key := fmt.Sprintf(hashKey, ut.String())
	hash := s.read(userBucket, key)

	return VerifyHash(string(hash), passphrase)
}

func (s *Store) ChangeUserPassword(ut UserToken, passphrase string) error {
	key := fmt.Sprintf(hashKey, ut.String())

	hash, err := GenerateHash(passphrase)
	if err != nil {
		return fmt.Errorf("could not Store.ChangeUserPassword: %v", err)
	}

	return s.write(userBucket, key, []byte(hash))
}

func (s *Store) GetFailedAuthCount(ut UserToken) (uint64, error) {
	key := fmt.Sprintf(failedKey, ut.String())

	i, err := s.readUint64(userBucket, key)
	if err != nil {
		return i, fmt.Errorf("could not Store.GetFailedAuthCount: %v", err)
	}

	return i, nil
}

func (s *Store) IncrementFailedAuthCount(ut UserToken) error {
	key := fmt.Sprintf(failedKey, ut.String())

	i, err := s.readUint64(userBucket, key)
	if err != nil {
		return fmt.Errorf("could not Store.IncrementFailedAuthCount: %v", err)
	}

	if i < maxFailCount {
		i = i + 1
	}

	return s.writeUint64(userBucket, key, i)
}

func (s *Store) ResetFailedAuthCount(ut UserToken) error {
	key := fmt.Sprintf(failedKey, ut.String())

	return s.writeUint64(userBucket, key, 0)
}
