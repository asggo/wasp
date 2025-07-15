package store

import (
	"fmt"
)

var (
	authFailedKey = "%s:failed"
	authHashKey   = "%s:hash"
	authTotpKey   = "%s:totp"
)

// ----------------------------------------------------------------------------
// Authentication Storage Methods
// ----------------------------------------------------------------------------
// Authenticate takes a passphrase and verifies it matches the user's original
// passphrase.
func (s *Store) AuthenticateUser(ut UserToken, passphrase, totp string) bool {
	count, err := getFailedAuthCount(ut)
	if err != nil {
		count = s.cfg.MaxAuthFailCount
	}

	// Sleep based on the failed auth count.
	time.Sleep(time.Duration(25*(1<<count)) * time.Millisecond)

	if !s.verifyHash(ut, string(hash), passphrase) {
		s.incrementFailedAuthCount(ut)
		return false
	}

	if !s.verifyTotp(ut, totp) {
		s.incrementFailedAuthCount(ut)
		return false
	}

	s.resetFailedAuthCount(ut)
	return true
}

func (s *Store) ChangeUserPassword(ut UserToken, passphrase string) error {
	key := fmt.Sprintf(hashKey, ut.String())

	hash, err := GenerateHash(passphrase)
	if err != nil {
		return fmt.Errorf("could not Store.ChangeUserPassword: %v", err)
	}

	return s.write(userBucket, key, []byte(hash))
}

func (s *Store) getFailedAuthCount(ut UserToken) (uint64, error) {
	key := fmt.Sprintf(authFailedKey, ut.String())

	i, err := s.readUint64(authBucket, key)
	if err != nil {
		return i, fmt.Errorf("could not Store.GetFailedAuthCount: %v", err)
	}

	return i, nil
}

func (s *Store) incrementFailedAuthCount(ut UserToken) error {
	key := fmt.Sprintf(authFailedKey, ut.String())

	i, err := s.readUint64(authBucket, key)
	if err != nil {
		return fmt.Errorf("could not Store.IncrementFailedAuthCount: %v", err)
	}

	if i < maxFailCount {
		i = i + 1
	}

	return s.writeUint64(authBucket, key, i)
}

func (s *Store) resetFailedAuthCount(ut UserToken) error {
	key := fmt.Sprintf(authFailedKey, ut.String())

	return s.writeUint64(authBucket, key, 0)
}
