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

func testSessionEqual(t *testing.T, s1, s2 Session) {
	if (s1.SessionId != s2.SessionId) || (s1.UserId != s2.UserId) || (s1.Expiration != s2.Expiration) {
		t.Fatal("Expected", s1, ", received", s2)
	}
}

func TestSession(t *testing.T) {
	fmt.Println(t.Name())

	u1 := NewUser(testUserAlias)
	s1, err := NewSession(u1.UserId, 5)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	bytes, err := s1.bytes()
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	s2, err := NewSessionFromBytes(bytes)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testSessionEqual(t, s1, s2)

	if s1.IsExpired() {
		t.Fatal("Expected unexpired session, received", s1)
	}

	s1.Expiration = 0
	time.Sleep(1)

	if !s1.IsExpired() {
		t.Fatal("Expected expired session, received", s1)
	}
}

func testStoreSession(t *testing.T) {
	fmt.Println(t.Name())

	u1 := NewUser(testUserAlias)
	s1, err := NewSession(u1.UserId, 5)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	db := newTestStore(t, testSessionDbPath)
	defer deleteTestStore(t, testSessionDbPath)

	// Create Session
	err = db.CreateSession(s1)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Get Session
	s2, err := db.GetSession(s1.SessionId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testSessionEqual(t, s1, s2)

	// Get Session From Request
	buf := new(bytes.Buffer)
	r, err := http.NewRequest("GET", "/", buf)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	r.AddCookie(&http.Cookie{Name: "sess", Value: s1.SessionId.String()})

	s3, err := NewSessionFromRequest(r, db)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	testSessionEqual(t, s1, s3)

	// Delete User
	err = db.DeleteSession(s1.SessionId)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	_, err = db.GetSession(s1.SessionId)
	if err == nil {
		t.Fatal("Expected error, received", nil)
	}

	db.Close()
}
