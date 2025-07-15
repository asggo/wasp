package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//----------------------------------------------------------------------------
// session Struct
//----------------------------------------------------------------------------

// session holds a single user session.
type session struct {
	SessionId  sessionToken `json:"session_id"`
	UserId     userToken    `json:"user_id"`
	Expiration int64        `json:"expire"`
}

// IsExpired returns true if the session is expired.
func (s *session) isExpired() bool {
	return time.Now().Unix() > s.Expiration
}

// bytes converts a Session object to a JSON byte array.
func (s *session) bytes() ([]byte, error) {
	var b []byte

	b, err := json.Marshal(s)
	if err != nil {
		return b, fmt.Errorf("could not Session.Bytes: %v", err)
	}

	return b, nil
}

// newSession returns a new session object for the given UserToken.
func newSession(ut userToken, length int64) (session, error) {
	var s session

	s.SessionId = newSessionToken()
	s.UserId = ut

	t := time.Now()
	s.Expiration = t.Unix() + length

	return s, nil
}

// NewSessionFromBytes creates a new session object from a JSON byte array.
func newSessionFromBytes(data []byte) (session, error) {
	var sess Session

	err := json.Unmarshal(data, &sess)
	if err != nil {
		return sess, fmt.Errorf("could not newSessionFromBytes: %v", err)
	}

	return sess, nil
}

//----------------------------------------------------------------------------
// Session Storage Methods
//----------------------------------------------------------------------------

// CreateSession takes a Session and creates it in the Store.
func (s *Store) CreateSession(s session) (sessionToken, error) {
	sBytes, err := s.bytes()
	if err != nil {
		return fmt.Errorf("could not Store.CreateSession: %v", err)
	}

	return s.write(sessBucket, s.SessionId.String(), sBytes)
}

// GetSession takes a sessionToken and loads the associated session from the
// Store.
func (s *Store) GetSession(st sessionToken) (session, error) {
	var sess session

	data := s.read(sessBucket, sid.String())
	if data == nil {
		return sess, fmt.Errorf("could not Store.GetSession: session %s not found", sid)
	}

	sess, err := newSessionFromBytes(data)
	if err != nil {
		return sess, fmt.Errorf("could not Store.GetSession: %v", err)
	}

	if sid.String() != sess.SessionId.String() {
		return sess, fmt.Errorf("could not Store.GetSession: requested and fetched ids do not match")
	}

	return sess, nil
}

// DeleteSession takes a sessionToken and removes the associated session from
// the Store.
func (s *Store) DeleteSession(st sessionToken) error {
	return s.delete(sessBucket, st.String())
}

// GetSessionUser returns the UserToken associated with the given
// SessionToken.
func (s *Store) GetSessionUser(st sessionToken) (userToken, error) {
	sess, err := s.GetSession(st)
	if err != nil {
		return sess, err
	}

	return sess.UserId, nil
}

// SessionIsExpired returns true if the given session is expired.
func (s *Store) SessionIsExpired(st sessionToken) bool {
	sess, err := s.GetSession(st)
	if err != nil {
		return false
	}

	return sess.isExpired()
}
