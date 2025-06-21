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

func testUserEqual(t *testing.T, u1, u2 User) {
	if (u1.UserId != u2.UserId) || (u1.Alias != u2.Alias) || (u1.Admin != u2.Admin) {
		t.Fatal("Expected", u1, ", received", u2)
	}
}

func TestUser(t *testing.T) {
	fmt.Println(t.Name())

	u1 := NewUser(testUserAlias)
	if u1.Admin {
		t.Fatal("Expected", false, ", received", u1.Admin)
	}

	bytes, err := u1.bytes()
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	u2, err := NewUserFromBytes(bytes)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testUserEqual(t, u1, u2)
}

func testStoreUser(t *testing.T) {
	fmt.Println(t.Name())

	u1 := NewUser(testUserAlias)
	s := newTestStore(t, testUserDbPath)
	defer deleteTestStore(t, testUserDbPath)

	// Create User
	err := s.CreateUser(u1, testUserPassphrase)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	err = s.CreateUser(u1, testUserPassphrase)
	if err == nil {
		t.Fatal("Expected error, received nil")
	}

	// User Exists
	if s.UserExists("nope") {
		t.Fatal("Expected user to not exist, but it does")
	}

	if !s.UserExists(testUserAlias) {
		t.Fatal("Expected user to exist, but it does not.")
	}

	// Get User
	u2, err := s.GetUser(u1.UserId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testUserEqual(t, u1, u2)

	// Get User by Alias
	u3, err := s.GetUserByAlias(testUserAlias)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testUserEqual(t, u1, u3)

	// Delete User
	err = s.DeleteUser(u1)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if s.UserExists(testUserAlias) {
		t.Fatal("Expected user to not exist, but it does.")
	}

	s.Close()
}
