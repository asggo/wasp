package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	userBucket = "user"
	sessBucket = "sess"
)

var (
	storeBuckets = [2]string{
		userBucket,
		sessBucket,
	}
)

// Store holds the bolt database
type Store struct {
	db *bolt.DB
}

// ----------------------------------------------------------------------------
// Helper Functions
// ----------------------------------------------------------------------------
func bytesToUint64(data []byte) (uint64, error) {
	var i uint64

	buf := bytes.NewBuffer(data)
	i, err := binary.ReadUvarint(buf)
	if err != nil {
		return i, err
	}

	return i, nil
}

func uint64ToBytes(i uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, i)

	return buf
}

//----------------------------------------------------------------------------
// Initialization Database
//----------------------------------------------------------------------------

// initialize configures the BBolt database for use as a Store.
func (s *Store) initialize() error {
	for _, bucket := range storeBuckets {
		err := s.createBucket(bucket)
		if err != nil {
			return err
		}
	}

	return nil
}

// createBucket creates a new bucket with the given name at the root of the
// database. An error is returned if the bucket cannot be created.
func (s *Store) createBucket(bucket string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return nil
	})
}

// ----------------------------------------------------------------------------
// Read, Write, Delete
// ----------------------------------------------------------------------------
// Write stores the given key/value pair in the given bucket.
func (s *Store) write(bucket, key string, value []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), value)
	})

	return err
}

// Read gets the value associated with the given key in the given bucket. If the
// key does not exist, Read returns nil.
func (s *Store) read(bucket, key string) []byte {
	var val []byte

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		val = b.Get([]byte(key))

		return nil
	})

	return val
}

// Delete removes a key/value pair from the given bucket. An error is returned
// if the key/value pair cannot be deleted.
func (s *Store) delete(bucket, key string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})

	return err
}

func (s *Store) readUint64(bucket, key string) (uint64, error) {
	var i uint64

	// Get existing value
	data := s.read(bucket, key)
	if data == nil {
		return i, nil
	} else {
		return bytesToUint64(data)
	}
}

func (s *Store) writeUint64(bucket, key string, i uint64) error {
	b := uint64ToBytes(i)

	return s.write(bucket, key, b)
}

//----------------------------------------------------------------------------
// Store Management
//----------------------------------------------------------------------------

// Backup creates a backup of the database to the given filename.
func (s *Store) Backup(filename string) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		file, e := os.Create(filename)
		if e != nil {
			return e
		}

		defer file.Close()

		_, e = tx.WriteTo(file)
		return e
	})

	return err
}

// Close closes the connection to the bbolt database.
func (s *Store) Close() error {
	return s.db.Close()
}

// NewStore creates a new Store object using a bbolt database located at the
// given filePath.
func NewStore(filePath string) (Store, error) {
	var s Store

	db, err := bolt.Open(filePath, 0640, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return s, fmt.Errorf("could not NewStore: %v", err)
	}

	s.db = db
	err = s.initialize()
	if err != nil {
		return s, fmt.Errorf("could not NewStore: %v", err)
	}

	return s, nil
}
