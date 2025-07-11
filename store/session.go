package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	sessionExpirationKey = "%s:expiration"
)

//----------------------------------------------------------------------------
// Session Storage Methods
//----------------------------------------------------------------------------

// CreateSession takes a Session and creates it in the Store.
func (s *Store) CreateSession(ut UserToken) (SessionToken, error) {
	st := NewSessionToken()
	now := time.Now().Unix()
	exp := now + s.cfg.SessionLength

	// Use a transaction to create the session in the database
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))

		// Associate sessionToken and the userToken
		err = b.Put([]byte(st.String()), []byte(ut.String()))
		if err != nil {
			return err
		}

		// Store the sessionToken expiration
		key := fmt.Sprintf(sessionExpirationKey, st.String())
		err = b.Put([]byte(key), uint64ToBytes(exp))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return st, fmt.Errorf("could not Store.CreateSession: %v", err)
	}

	return st, nil
}

// DeleteSession takes a SessionToken and removes the associated session from
// the Store.
func (s *Store) DeleteSession(st SessionToken) error {
	// Use a transaction to delete the session from the database
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))

		// Delete the sessionToken expiration
		key := fmt.Sprintf(sessionExpirationKey, st.String())
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Delete the sessionToken
		err = b.Delete([]byte(st.String()))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not Store.DeleteSession: %v", err)
	}

	return nil
}

// GetSessionUser returns the UserToken associated with the given
// SessionToken.
func (s *Store) GetSessionUser(st SessionToken) (UserToken, error) {
	var ut UserToken

	data := s.read(sessBucket, sid.String())
	if data == nil {
		return ut, fmt.Errorf("could not Store.GetSessionUser: session %s not found", sid)
	}

	ut, err := parseUserToken(string(data))
	if err != nil {
		return ut, fmt.Errorf("could not Store.GetSessionUser: %v", err)
	}

	return ut, nil
}

// SessionIsExpired returns true if the given session is expired.
func (s *Store) SessionIsExpired(st SessionToken) bool {
	key := fmt.Sprintf(sessionExpirationKey, st.String())
	current := time.Now().Unix()

	expiration, err := s.readUint64(sessionBucket, key)
	if err != nil {
		return true
	}

	return current > expiration
}
