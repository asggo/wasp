package store

import (
	"fmt"
	"net/http"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	sessionUserKey   = "%s:userid"
	sessionExpireKey = "%s:expire"
)

//----------------------------------------------------------------------------
// Session Struct
//----------------------------------------------------------------------------

// Session holds a single user session.
type Session struct {
	SessionId  SessionToken `json:"session_id"`
	UserId     UserToken    `json:"user_id"`
	Expiration uint64       `json:"expire"`
}

// IsExpired returns true if the session is expired.
func (s *Session) IsExpired() bool {
	t := time.Now()

	return t.Unix() > int64(s.Expiration)
}

// NewSession returns a new Session object for the given User.
func NewSession(uid UserToken, length int64) (Session, error) {
	var s Session

	s.SessionId = NewSessionToken()
	s.UserId = uid

	t := time.Now()
	s.Expiration = uint64(t.Unix() + length)

	return s, nil
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

	sess, err = s.ReadSession(sessId)
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
	err := s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessBucket))

		// Write each of the session elements to the Store
		key := fmt.Sprintf(sessionUserKey, sess.SessionId.String())
		err := b.Put([]byte(key), []byte(sess.UserId.String()))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(sessionExpireKey, sess.SessionId.String())
		err = b.Put([]byte(key), uint64ToBytes(sess.Expiration))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not store.CreateSession: %v", err)
	}

	return nil
}

// ReadSession takes a SessionToken and returns the Session associated with it.
func (s *Store) ReadSession(sid SessionToken) (Session, error) {
	var sess Session
	var userId []byte
	var expire []byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessBucket))

		// Read each of the session elements from the Store
		key := fmt.Sprintf(sessionUserKey, sid.String())
		userId = b.Get([]byte(key))
		if userId == nil {
			return fmt.Errorf("no sessionUserKey")
		}

		key = fmt.Sprintf(sessionExpireKey, sid.String())
		expire = b.Get([]byte(key))
		if expire == nil {
			return fmt.Errorf("no sessionExpireKey")
		}

		return nil
	})

	if err != nil {
		return sess, fmt.Errorf("could not store.ReadSession: %v", err)
	}

	token, err := parseUserToken(string(userId))
	if err != nil {
		return sess, fmt.Errorf("could not store.ReadSession: %v", err)
	}

	n, err := bytesToUint64(expire)
	if err != nil {
		return sess, fmt.Errorf("could not store.ReadSession: %v", err)
	}

	sess.SessionId = sid
	sess.UserId = token
	sess.Expiration = n

	return sess, nil
}

// DeleteSession takes a SessionToken and removes the associated session from
// the Store.
func (s *Store) DeleteSession(sid SessionToken) error {
	err := s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessBucket))

		// Delete each of the session elements from the Store
		key := fmt.Sprintf(sessionUserKey, sid.String())
		err := b.Delete([]byte(key))
		if err != nil {
			return err
		}

		key = fmt.Sprintf(sessionExpireKey, sid.String())
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not store.DeleteSession: %v", err)
	}

	return nil
}
