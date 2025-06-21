package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//----------------------------------------------------------------------------
// Session Struct
//----------------------------------------------------------------------------

// Session holds a single user session.
type Session struct {
	SessionId  SessionToken `json:"session_id"`
	UserId     UserToken    `json:"user_id"`
	Expiration int64        `json:"expire"`
}

// IsExpired returns true if the session is expired.
func (s *Session) IsExpired() bool {
	t := time.Now()

	return t.Unix() > s.Expiration
}

// bytes converts a Session object to a JSON byte array.
func (s *Session) bytes() ([]byte, error) {
	var b []byte

	b, err := json.Marshal(s)
	if err != nil {
		return b, fmt.Errorf("could not Session.Bytes: %v", err)
	}

	return b, nil
}

// NewSession returns a new Session object for the given User.
func NewSession(uid UserToken, length int64) (Session, error) {
	var s Session

	s.SessionId = NewSessionToken()
	s.UserId = uid

	t := time.Now()
	s.Expiration = t.Unix() + length

	return s, nil
}

// NewSessionFromBytes creates a new Session object from a JSON byte array.
func NewSessionFromBytes(data []byte) (Session, error) {
	var sess Session

	err := json.Unmarshal(data, &sess)
	if err != nil {
		return sess, fmt.Errorf("could not NewSessionFromStore: %v", err)
	}

	return sess, nil
}

// NewSessionFromRequest loads a Session from the Store based on the session
// cookie in the given HTTP request.
func NewSessionFromRequest(r *http.Request, s *Store) (Session, error) {
	var sess Session

	sessCookie, err := r.Cookie("sess")
	if err != nil {
		return sess, fmt.Errorf("could not NewSessionFromRequest: %v", err)
	}

	sessId, err := parseSessionToken(sessCookie.Value)
	if err != nil {
		return sess, fmt.Errorf("could not NewSessionFromRequest: %v", err)
	}

	sess, err = s.GetSession(sessId)
	if err != nil {
		return sess, fmt.Errorf("could not NewSessionFromRequest: %v", err)
	}

	return sess, nil
}

//----------------------------------------------------------------------------
// Session Storage Methods
//----------------------------------------------------------------------------

// CreateSession takes a Session and creates it in the Store.
func (s *Store) CreateSession(sess Session) error {
	sessionBytes, err := sess.bytes()
	if err != nil {
		return fmt.Errorf("could not Store.CreateSession: %v", err)
	}

	return s.write(sessBucket, sess.SessionId.String(), sessionBytes)
}

// DeleteSession takes a SessionToken and removes the associated session from
// the Store.
func (s *Store) DeleteSession(sid SessionToken) error {
	return s.delete(sessBucket, sid.String())
}

// GetSession takes a SessionToken and returns the Session associated with it.
func (s *Store) GetSession(sid SessionToken) (Session, error) {
	var sess Session

	data := s.read(sessBucket, sid.String())
	if data == nil {
		return sess, fmt.Errorf("could not Store.GetSession: session %s not found", sid)
	}

	sess, err := NewSessionFromBytes(data)
	if err != nil {
		return sess, fmt.Errorf("could not Store.GetSession: %v", err)
	}

	if sid.String() != sess.SessionId.String() {
		return sess, fmt.Errorf("could not Store.GetSession: requested and fetched ids do not match")
	}

	return sess, nil
}
