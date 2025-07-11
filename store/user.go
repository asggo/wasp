package store

import (
	"encoding/json"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/text/unicode/norm"
)

var (
	userFailedKey = "%s:failed"
	userHashKey   = "%s:hash"
	userTotpKey   = "%s:totp"
	userAdminKey  = "%s:admin"
)

//----------------------------------------------------------------------------
// User Storage Methods
//----------------------------------------------------------------------------

// CreateUser takes a UserToken, alias, password, and an admin flag and
// creates a user in the Store.
func (s *Store) CreateUser(alias, pwd string, adm bool) error {
	// Verify the alias does not already exist
	data := s.read(userBucket, alias)
	if data != nil {
		return fmt.Errorf("could not Store.CreateUser: alias %s exists", alias)
	}

	// Create a UserToken for the user.
	ut := NewUserToken()

	// Create a TotpToken for the user.
	secret := NewTotpToken()

	// Generate the user's password hash.
	ah, err := newArgonHash(s.cfg.ArgonTime, s.cfg.ArgonMemory, s.cfg.ArgonThreads)
	if err != nil {
		return fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	hash := ah.derive(pwd)

	// Use a transaction to create the user account in the database
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Associate alias and user id
		err = b.Put([]byte(alias), []byte(ut.String()))
		if err != nil {
			return err
		}

		// Store the user's password hash
		key := fmt.Sprintf(userHashKey, ut.String())
		err = b.Put([]byte(key), []byte(hash))
		if err != nil {
			return err
		}

		// Store the user's failed authentication count
		key = fmt.Sprintf(userFailedKey, ut.String())
		err = b.Put([]byte(key), uint64ToBytes(0))
		if err != nil {
			return err
		}

		// Store the user's totp secret
		key = fmt.Sprintf(userTotpKey, ut.String())
		err = b.Put([]byte(key), []byte(secret.String()))
		if err != nil {
			return err
		}

		// Set the user admin flag
		key = fmt.Sprintf(userAdminKey, ut.String())
		err = b.Put([]byte(key), []byte(fmt.Sprintf("%t", adm)))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	return nil
}

// DeleteUser takes an alias and removes all the keys associated with the
// user from the Store.
func (s *Store) DeleteUser(alias string) error {
	ut := s.read(userBucket, alias)
	if data == nil {
		return fmt.Errorf("could not Store.DeleteUser: alias %s does not exist", alias)
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Remove the admin flag
		key := fmt.Sprintf(userAdminKey, string(ut))
		err := b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the Totp secret
		key = fmt.Sprintf(userTotpKey, string(ut))
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the password hash
		key = fmt.Sprintf(userHashKey, string(ut))
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the failed authentication count
		key = fmt.Sprintf(userFailedKey, string(ut))
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the alias
		err = b.Delete([]byte(alias))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not Store.DeleteUser: %v", err)
	}

	return nil
}

// GetUserToken takes an alias and returns the UserToken associated with it.
func (s *Store) GetUserToken(alias string) (UserToken, error) {
	var ut UserToken

	data := s.read(userBucket, alias)
	if data == nil {
		return ut, fmt.Errorf("could not Store.GetUserToken: alias %s not found", alias)
	}

	ut, err := parseUserToken(string(data))
	if err != nil {
		return ut, fmt.Errorf("could not Store.GetUserByAlias: %s %v", alias, err)
	}

	return ut, nil
}

// UserExists returns true if the given alias is already registered.
func (s *Store) UserExists(alias string) bool {
	data := s.read(userBucket, un)

	// A nil result means the user does not exist.
	return data != nil
}
