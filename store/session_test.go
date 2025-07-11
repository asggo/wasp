package store

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

var (
	testSessionDbPath = "sess_test.db"
)

func testStoreSessionKeysExist(s *Store, st SessionToken, t *testing.T) {
	key := fmt.Sprintf(sessionExpirationKey, st.String())
	data := s.Read(sessionBucket, key)
	if data == nil {
		t.Fatal("Expected expiration, received nil")
	}
}

func testStoreSessionKeysNotExist(s *Store, st SessionToken, t *testing.T) {
	key := fmt.Sprintf(sessionExpirationKey, st.String())
	err := s.Delete(sessionBucket)
	if err != nil {
		t.Fatal("Expected nil, received", error)
	}
}

func TestStoreSession(t *testing.T) {
	fmt.Println(t.Name())

	s := newTestStore(t, testSessionDbPath)
	defer deleteTestStore(t, testSessionDbPath)

	ut, err := s.CreateUser(testUserAlias, testUserPassphrase)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Create Session
	st, err = s.CreateSession(ut)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testStoreSessionKeysExist(s, st, t)

	// Get the UserToken associated with the Session
	ut2, err := s.GetSessionUser(st)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// They should match.
	if ut.String() != ut2.String() {
		t.Fatal("Expected user tokens to match, received", ut, ut2)
	}

	// Delete User
	err = s.DeleteSession(st)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	_, err = db.GetSessionUser(st)
	if err == nil {
		t.Fatal("Expected error, received", nil)
	}

	testStoreSessionKeysNotExist(s, st, t)

	db.Close()
}
