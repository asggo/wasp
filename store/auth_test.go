package store

import (
	"fmt"
	"testing"
)

var (
	testAuthGoodPassword = "testauthgoodpassword"
	testAuthBadPassword  = "testauthbadpassword"
	testAuthDbPath       = "auth_test.db"
)

func testStoreAuth(t *testing.T) {
	fmt.Println(t.Name())

	u1 := NewUser(testUserAlias)
	db := newTestStore(t, testAuthDbPath)
	defer deleteTestStore(t, testAuthDbPath)

	err := db.CreateUser(u1, testAuthGoodPassword)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Verify we can authenticate with the right password.
	if !db.AuthenticateUser(u1.UserId, testAuthGoodPassword) {
		t.Fatal("Expected password to match:", testAuthGoodPassword)
	}

	// Verify we cannot authenticate with the wrong password.
	if db.AuthenticateUser(u1.UserId, testAuthBadPassword) {
		t.Fatal("Expected password to not match:", testAuthBadPassword)
	}

	// Change password
	err = db.ChangeUserPassword(u1.UserId, testAuthBadPassword)
	if err != nil {
		t.Fatal("Expected", nil, ", recieved", err)
	}

	// Verify we can authenticate with the new password.
	if !db.AuthenticateUser(u1.UserId, testAuthBadPassword) {
		t.Fatal("Expected password to match:", testAuthBadPassword)
	}

	// Get failed auth count
	count, err := db.GetFailedAuthCount(u1.UserId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if count != 0 {
		t.Fatal("Expected", 0, ", received", count)
	}

	// Increment the failed auth count
	err = db.IncrementFailedAuthCount(u1.UserId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	count, err = db.GetFailedAuthCount(u1.UserId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if count != 1 {
		t.Fatal("Expected", 1, ", received", count)
	}

	// Reset the failed auth count
	err = db.ResetFailedAuthCount(u1.UserId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	count, _ = db.GetFailedAuthCount(u1.UserId)
	if count != 0 {
		t.Fatal("Expected", 0, ", received", count)
	}

	db.Close()
}
