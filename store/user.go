package store

import (
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/text/unicode/norm"
)

var (
	userAliasKey  = "%s:alias"
	userAdminKey  = "%s:admin"
	userHashKey   = "%s:hash"
	userFailedKey = "%s:failed"
)

//----------------------------------------------------------------------------
// User Struct
//----------------------------------------------------------------------------

// User holds a single user account.
type User struct {
	UserId       UserToken
	Alias        string
	Admin        bool
	PasswordHash string
	FailedCount  uint64
}

// NewUser creates a new User object using the given alias and passphrase. The
// passphrase is hashed and stored in the User object.
func NewUser(alias, passphrase string) User {
	var u User

	u.UserId = NewUserToken()
	u.Alias = strings.ToLower(norm.NFKD.String(alias))
	u.Admin = false
	u.PasswordHash = GenerateHash(passphrase)
	u.FailedCount = 0

	return u
}

//----------------------------------------------------------------------------
// User Storage Methods
//----------------------------------------------------------------------------

// CreateUser takes a User and creates it in the Store. First it verifies the
// User alias does not already exist, then it uses a transaction to associate
// the alias to the user id and to write each element of the User object to
// the store.
func (s *Store) CreateUser(u User) error {
	// Verify the alias does not already exist
	data := s.read(userBucket, u.Alias)
	if data != nil {
		return fmt.Errorf("could not Store.CreateUser: alias %s exists", u.Alias)
	}

	// Create our User in the Store
	err := s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Associate alias and user id
		err := b.Put([]byte(u.Alias), []byte(u.UserId.String()))
		if err != nil {
			return err
		}

		// Write each of the user elements to the Store
		key := fmt.Sprintf(userAliasKey, u.UserId.String())
		err = b.Put([]byte(key), []byte(u.Alias))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userAdminKey, u.UserId.String())
		err = b.Put([]byte(key), boolToBytes(u.Admin))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userHashKey, u.UserId.String())
		err = b.Put([]byte(key), []byte(u.PasswordHash))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userFailedKey, u.UserId.String())
		err = b.Put([]byte(key), uint64ToBytes(u.FailedCount))
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

// ReadUser takes a UserToken and returns the user associated with it.
func (s *Store) ReadUser(uid UserToken) (User, error) {
	var user User

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Read each of the user elements from the Store
		key := fmt.Sprintf(userAliasKey, uid.String())
		alias := b.Get([]byte(key))
		if alias == nil {
			return fmt.Errorf("no userAliasKey")
		}

		key = fmt.Sprintf(userAdminKey, uid.String())
		admin := b.Get([]byte(key))
		if admin == nil {
			return fmt.Errorf("no userAdminKey")
		}

		key = fmt.Sprintf(userHashKey, uid.String())
		hash := b.Get([]byte(key))
		if hash == nil {
			return fmt.Errorf("no userHashKey")
		}

		key = fmt.Sprintf(userFailedKey, uid.String())
		failed := b.Get([]byte(key))
		if failed == nil {
			return fmt.Errorf("no userFailedKey")
		}

		user.Alias = string(alias)
		user.Admin = bytesToBool(admin)
		user.PasswordHash = string(hash)

		count, err := bytesToUint64(failed)
		if err != nil {
			return err
		}

		user.FailedCount = count

		return nil
	})

	if err != nil {
		return user, fmt.Errorf("could not Store.GetUser: %v", err)
	}

	user.UserId = uid

	return user, nil
}

// DeleteUser takes a User and removes it from the Store. A transaction is
// used to delete all the keys associated with the user.
func (s *Store) DeleteUser(u User) error {
	err := s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Delete the Alias from the Store
		err := b.Delete([]byte(u.Alias))
		if err != nil {
			return err
		}

		// Delete each of the user elements from the Store
		key := fmt.Sprintf(userAliasKey, u.UserId.String())
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userAdminKey, u.UserId.String())
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userHashKey, u.UserId.String())
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(userFailedKey, u.UserId.String())
		err = b.Delete([]byte(key))
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

// GetUserByAlias takes an alias and returns the User associated with it.
func (s *Store) GetUserByAlias(alias string) (User, error) {
	var user User

	data := s.read(userBucket, alias)
	if data == nil {
		return user, fmt.Errorf("could not Store.GetUserByAlias: alias %s not found", alias)
	}

	token, err := parseUserToken(string(data))
	if err != nil {
		return user, fmt.Errorf("could not Store.GetUserByAlias: %s %v", alias, err)
	}

	return s.ReadUser(token)
}

// UserExists returns true if the given alias is already registered.
func (s *Store) UserExists(alias string) bool {
	data := s.read(userBucket, alias)

	// A nil result means the user does not exist.
	return data != nil
}
