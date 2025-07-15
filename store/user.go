package store

import (
	"encoding/json"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/text/unicode/norm"
)

//----------------------------------------------------------------------------
// User Struct
//----------------------------------------------------------------------------

// user holds a single user account.
type user struct {
	UserId userToken `json:"user_id"`
	Alias  string    `json:"alias"`
	Admin  bool      `json:"admin"`
}

// bytes renders a user object as a JSON byte array.
func (u *user) bytes() ([]byte, error) {
	var b []byte

	b, err := json.Marshal(u)
	if err != nil {
		return b, fmt.Errorf("could not user.bytes: %v", err)
	}

	return b, nil
}

// newUser creates a new User object using the given alias and passphrase. The
// passphrase is hashed and stored in the User object.
func newUser(alias string, admin bool) user {
	var u user

	u.UserId = newUserToken()
	u.Alias = strings.ToLower(norm.NFKD.String(alias))
	u.Admin = admin

	return u
}

// newUserFromBytes creates a new user object from a JSON byte array.
func newUserFromBytes(data []byte) (user, error) {
	var u user

	err := json.Unmarshal(data, &u)
	if err != nil {
		fmt.Println(string(data))
		return user, fmt.Errorf("could not NewUserFromBytes: %v", err)
	}

	return u, nil
}

//----------------------------------------------------------------------------
// User Storage Methods
//----------------------------------------------------------------------------

// CreateUser takes a UserToken, alias, password, and an admin flag and
// creates a user in the Store.
func (s *Store) CreateUser(alias, pwd string, admin bool) (userToken, error) {
	var ut userToken

	// Verify the alias does not already exist
	data := s.read(userBucket, alias)
	if data != nil {
		return ut, fmt.Errorf("could not Store.CreateUser: alias %s exists", alias)
	}

	// Create a user.
	user := newUser(alias, admin)
	ut = user.UserId

	// Create a TotpToken for the user.
	totp := newTotp()

	// Create the user's passwordHash
	ph, err := newPasswordHash(s.cfg.ArgonTime, s.cfg.ArgonMemory, s.cfg.ArgonThreads, pwd)
	if err != nil {
		return ut, fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	// Use a transaction to create the user account in the database
	err = s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(userBucket))
		ab := tx.Bucket([]byte(authBucket))
		uid := user.UserId.String()

		// Associate alias and user id in the userBucket
		err = ub.Put([]byte(alias), []byte(uid))
		if err != nil {
			return err
		}

		// Get the user object bytes and store them in the userBucket
		data, err := user.bytes()
		if err != nil {
			return err
		}

		err = ub.Put([]byte(uid), data)
		if err != nil {
			return err
		}

		// Get the passwordHash bytes and store them in the authBucket
		data, err = ph.bytes()
		if err != nil {
			return err
		}

		key := fmt.Sprintf(authHashKey, uid)
		err = ab.Put([]byte(key), data)
		if err != nil {
			return err
		}

		// Get the totp object bytes and store them in the authBucket
		data, err = totp.bytes()
		if err != nil {
			return err
		}

		key = fmt.Sprintf(authTotpKey, uid)
		err = ab.Put([]byte(key), data)
		if err != nil {
			return err
		}

		// Store the user's failed authentication count in the authBucket
		key = fmt.Sprintf(authFailedKey, uid)
		err = ab.Put([]byte(key), uint64ToBytes(0))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ut, fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	return ut, nil
}

// DeleteUser takes an alias and removes all the keys associated with the user
// from the Store.
func (s *Store) DeleteUser(alias string) error {
	ut := s.read(userBucket, alias)
	if ut == nil {
		return fmt.Errorf("could not Store.DeleteUser: alias %s does not exist", alias)
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(userBucket))
		ab := tx.Bucket([]byte(authBucket))

		// Remove the failed auth count from the authBucket
		key = fmt.Sprintf(authFailedKey, string(ut))
		err = ab.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the totp object from the authBucket
		key = fmt.Sprintf(authTotpKey, string(ut))
		err = ab.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the passwordHash object from the authBucket
		key = fmt.Sprintf(authHashKey, string(ut))
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove the user object from the userBucket
		err := b.Delete(ut)
		if err != nil {
			return err
		}

		// Remove the alias association from the userBucket
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
