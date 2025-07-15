package store

import (
	"fmt"
	"testing"
)

var (
	testAuthGoodPassword = "testauthgoodpassword"
	testAuthBadPassword  = "testauthbadpassword"
	testAuthBadCode      = "000000"
	testAuthDbPath       = "auth_test.db"
)

func testFailedAuthCount(s *Store, ut UserToken, t *testing.T) {
	fmt.Println(t.Name())

	// Get failed auth count
	count, err := s.getFailedAuthCount(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if count != 0 {
		t.Fatal("Expected", 0, ", received", count)
	}

	// Increment the failed auth count
	err = s.incrementFailedAuthCount(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	count, err = s.getFailedAuthCount(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if count != 1 {
		t.Fatal("Expected", 1, ", received", count)
	}

	// Reset the failed auth count
	err = s.resetFailedAuthCount(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	count, _ = s.getFailedAuthCount(ut)
	if count != 0 {
		t.Fatal("Expected", 0, ", received", count)
	}
}

func testStoreAuth(t *testing.T) {
	fmt.Println(t.Name())

	s := newTestStore(t, testAuthDbPath)
	defer deleteTestStore(t, testAuthDbPath)

	// Create a new user account
	ut, err := s.CreateUser(testUserAlias, testAuthGoodPassword)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Test the failed auth count functions.
	testFailedAuthCount(s, ut, t)

	// Get the users TotpToken so we can generate an auth code.
	tt, err := s.GetUserTotpToken(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Generate an auth code
	code := generateSha256Totp(
		tt[:],
		time.Now().Unix(),
		s.cfg.TotpLength,
		s.cfg.TotpStart,
		s.cfg.TotpStep,
	)

	// Verify we cannot authenticate with the wrong password.
	if s.AuthenticateUser(testUserAlias, testAuthBadPassword, code) {
		t.Fatal("Expected password to not match:", testAuthBadPassword)
	}

	// Verify we cannot authenticate with the wrong code.
	if s.AuthenticateUser(testUserAlias, testAuthGoodPassword, testAuthBadCode) {
		t.Fatal("Expected TOTP codes to not match:", testAuthBadCode)
	}

	// FailedAuthCount should be 2 at this point
	count, _ := s.getFailedAuthCount(ut)
	if count != 2 {
		t.Fatal("Expected 2 failed authentication attempts, received", count)
	}

	// Verify we can authenticate with the right password and code.
	if !s.AuthenticateUser(testUserAlias, testAuthGoodPassword, code) {
		t.Fatal("Expected successful login:", testAuthGoodPassword, code)
	}

	// The count should be 0 after the successful login
	count, _ = s.getFailedAuthCount(ut)
	if count != 0 {
		t.Fatal("Expected 0 failed authentication attempts, received", count)
	}

	// Verify we cannot authenticate with the right password and code a second
	// time.
	if s.AuthenticateUser(testUserAlias, testAuthGoodPassword, code) {
		t.Fatal("Expected invalid TOTP code:", code)
	}

	// Change password
	err = s.ChangeUserPassword(ut, testAuthBadPassword)
	if err != nil {
		t.Fatal("Expected", nil, ", recieved", err)
	}

	// Generate an auth code
	time.Sleep(s.cfg.TotpStep)
	code = generateSha256Totp(
		tt[:],
		time.Now().Unix(),
		s.cfg.TotpLength,
		s.cfg.TotpStart,
		s.cfg.TotpStep,
	)

	// Verify we can authenticate with the new password and code.
	if !s.AuthenticateUser(testUserAlias, testAuthBadPassword, code) {
		t.Fatal("Expected successful login:", testAuthBadPassword, code)
	}

	// Delete the user's TotpToken
	err = s.DeleteUserTotpToken(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	time.Spleep(s.cfg.TotpStep)
	code = generateSha256Totp(
		tt[:],
		time.Now().Unix(),
		s.cfg.TotpLength,
		s.cfg.TotpStart,
		s.cfg.TotpStep,
	)

	if db.AuthenticateUser(testUserAlias, testAuthBadPassword, code) {
		t.Fatal("Expected failed login", testAuthBadPassword, code)
	}

	db.Close()
}
