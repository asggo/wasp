package store

import (
	"fmt"
	"testing"
)

var (
	testUserAlias      = "alias"
	testUserPassphrase = "aliaspassword123"
	testUserDbPath     = "user_test.db"
)

func testStoreUserKeysExist(s *Store, t *testing.T) {
	key := fmt.Sprintf(userFailedKey, string(ut))
	data := s.Read(userBucket, key)
	if data == nil {
		t.Fatal("Expected failed count, received nil")
	}

	key = fmt.Sprintf(userHashKey, string(ut))
	data = s.Read(userBucket, key)
	if data == nil {
		t.Fatal("Expected password hash, received nil")
	}

	key = fmt.Sprintf(userTotpKey, string(ut))
	data = s.Read(userBucket, key)
	if data == nil {
		t.Fatal("Expected TotpKey, received nil")
	}

	key = fmt.Sprintf(userAdminKey, string(ut))
	data = s.Read(userBucket, key)
	if data == nil {
		t.Fatal("Expected Admin flag, received nil")
	}
}

func testStoreUserKeysNotExist(s *Store, t *testing.T) {
	key := fmt.Sprintf(userFailedKey, string(ut))
	data := s.Read(userBucket, key)
	if data != nil {
		t.Fatal("Expected nil, received", data)
	}

	key = fmt.Sprintf(userHashKey, string(ut))
	data = s.Read(userBucket, key)
	if data != nil {
		t.Fatal("Expected nil, received", data)
	}

	key = fmt.Sprintf(userTotpKey, string(ut))
	data = s.Read(userBucket, key)
	if data != nil {
		t.Fatal("Expected nil, received", data)
	}

	key = fmt.Sprintf(userAdminKey, string(ut))
	data = s.Read(userBucket, key)
	if data != nil {
		t.Fatal("Expected nil, received", data)
	}
}

func testStoreUserExists(s *Store, t *testing.T) {
	if !s.UserExists(testUserAlias) {
		t.Fatal("Expected user to exist, but it does not")
	}

	ut, err = GetUserToken(testUserAlias)
	if err != nil {
		t.Fatal("Expected no error, received", err)
	}

}

func testStoreUserNotExists(s *Store, t *testing.T) {
	fmt.Println(t.Name())

	if s.UserExists(testUserAlias) {
		t.Fatal("Expected user to not exist, but it does")
	}

	ut, err := GetUserToken(testUserAlias)
	if err == nil {
		t.Fatal("Expected an error, received nil")
	}
}

func TestStoreUser(t *testing.T) {
	fmt.Println(t.Name())

	s := newTestStore(t, testUserDbPath)
	defer deleteTestStore(t, testUserDbPath)

	testStoreUserNotExists(s, t)
	testStoreUserKeysNotExist(s, t)

	// Create the user
	err := s.CreateUser(testUserAlias, testUserPassphrase, false)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// The user is created so CreateUser should return an error.
	err = s.CreateUser(testUserAlias, testUserPassphrase, false)
	if err == nil {
		t.Fatal("Expected error, received nil")
	}

	testStoreUserExists(s, t)
	testStoreUserKeysExist(s, t)

	// Delete User
	err = s.DeleteUser(testUserAlias)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testStoreUserNotExists(s, t)
	testStoreUserKeysNotExist(s, t)

	s.Close()
}
