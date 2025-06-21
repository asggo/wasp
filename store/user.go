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

// User holds a single user account.
type User struct {
	UserId UserToken `json:"user_id"`
	Alias  string    `json:"alias"`
	Admin  bool      `json:"admin"`
}

// bytes renders a User object as a JSON byte array.
func (u *User) bytes() ([]byte, error) {
	var b []byte

	b, err := json.Marshal(u)
	if err != nil {
		return b, fmt.Errorf("could not User.bytes: %v", err)
	}

	return b, nil
}

// NewUser creates a new User object using the given alias and passphrase. The
// passphrase is hashed and stored in the User object.
func NewUser(alias string) User {
	var u User

	u.UserId = NewUserToken()
	u.Alias = strings.ToLower(norm.NFKD.String(alias))
	u.Admin = false

	return u
}

// NewUserFromBytes creates a new User object from a JSON byte array.
func NewUserFromBytes(data []byte) (User, error) {
	var user User

	err := json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println(string(data))
		return user, fmt.Errorf("could not NewUserFromBytes: %v", err)
	}

	return user, nil
}

//----------------------------------------------------------------------------
// User Storage Methods
//----------------------------------------------------------------------------

// CreateUser takes a User and creates it in the Store. A transaction is used
// to create two keys, one to relate the alias to the user id and the other to
// relate the user id to the User bytes.
func (s *Store) CreateUser(u User, passphrase string) error {
	userBytes, err := u.bytes()
	if err != nil {
		return fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	hash, err := GenerateHash(passphrase)
	if err != nil {
		return fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	// Verify the alias does not already exist
	data := s.read(userBucket, u.Alias)
	if data != nil {
		return fmt.Errorf("could not Store.CreateUser: alias %s exists", u.Alias)
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		// Associate alias and user id
		err = b.Put([]byte(u.Alias), []byte(u.UserId.String()))
		if err != nil {
			return err
		}

		// Associate user id and User bytes
		err = b.Put([]byte(u.UserId.String()), userBytes)
		if err != nil {
			return err
		}

		// Store the user's password hash
		key := fmt.Sprintf(hashKey, u.UserId.String())
		err = b.Put([]byte(key), []byte(hash))
		if err != nil {
			return err
		}

		// Store the user's failed authentication count
		key = fmt.Sprintf(failedKey, u.UserId.String())
		err = b.Put([]byte(key), uint64ToBytes(0))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not Store.CreateUser: %v", err)
	}

	// Add the user's passwordHash to the store.

	return nil
}

// // SaveUser takes a User and updates it in the Store.
// func (s *Store) SaveUser(u User) error {
// 	userBytes, err := u.bytes()
// 	if err != nil {
// 		return fmt.Errorf("could not Store.SaveUser: %v", err)
// 	}

// 	return s.write(userBucket, u.UserId.String(), userBytes)
// }

// DeleteUser takes a User and removes it from the Store. A transaction is
// used to delete both keys associated with the user.
func (s *Store) DeleteUser(u User) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		err := b.Delete([]byte(u.UserId.String()))
		if err != nil {
			return err
		}

		err = b.Delete([]byte(u.Alias))
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

// GetUser takes a UserToken and returns the user associated with it.
func (s *Store) GetUser(uid UserToken) (User, error) {
	var user User

	data := s.read(userBucket, uid.String())
	if data == nil {
		return user, fmt.Errorf("could not Store.GetUser: user %s not found", uid)
	}

	user, err := NewUserFromBytes(data)
	if err != nil {
		return user, fmt.Errorf("could not Store.GetUser: %v", err)
	}

	if uid.String() != user.UserId.String() {
		return user, fmt.Errorf("could not Store.GetUser: requested and fetched ids do not match")
	}

	return user, nil
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

	return s.GetUser(token)
}

// UserExists returns true if the given user is already registered.
func (s *Store) UserExists(un string) bool {
	data := s.read(userBucket, un)

	// A nil result means the user does not exist.
	return data != nil
}
