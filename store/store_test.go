package store

import (
	"fmt"
	"os"
	"testing"
)

var (
	testCoreKeyName      = "core_key"
	testCoreVal1         = "core_value_1"
	testCoreVal2         = "core_value_2"
	testCoreIntKey       = "integer"
	testCoreInt          = uint64(1024)
	testCoreDbPath       = "core_test.db"
	testCoreBackupDbPath = "core_test_bu.db"
)

func TestStore(t *testing.T) {
	t.Run("Test Store Core", testStoreCore)
	t.Run("Test Store Backup", testStoreBackup)
	t.Run("Test Store Auth", testStoreAuth)
	t.Run("Test Store User", testStoreUser)
	t.Run("Test Store Session", testStoreSession)
}

func newTestStore(t *testing.T, path string) *Store {
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("could not newTestStore: %v", err)
	}

	return &s
}

func deleteTestStore(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil {
		t.Fatalf("could not deleteTestStore: %v", err)
	}
}

func testStoreBucketWRD(t *testing.T, db *Store, bucket string) {
	// Write to the bucket
	err := db.write(bucket, testCoreKeyName, []byte(testCoreVal1))
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Read from the bucket
	val := db.read(bucket, testCoreKeyName)
	if string(val) != testCoreVal1 {
		t.Fatal("Expected", testCoreVal1, ", received", val)
	}

	// Update the key
	db.write(bucket, testCoreKeyName, []byte(testCoreVal2))
	val = db.read(bucket, testCoreKeyName)
	if string(val) != testCoreVal2 {
		t.Fatal("Expected", testCoreVal2, ", received", val)
	}

	// Delete the key
	err = db.delete(bucket, testCoreKeyName)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	// Read the key
	val = db.read(bucket, testCoreKeyName)
	if val != nil {
		t.Fatal("Expected nil value, received", val)
	}
}

func testStoreCore(t *testing.T) {
	fmt.Println(t.Name())

	db := newTestStore(t, testCoreDbPath)
	defer deleteTestStore(t, testCoreDbPath)

	// Confirm we can write to the two buckets we initialized.
	for _, b := range storeBuckets {
		testStoreBucketWRD(t, db, b)
	}

	err := db.writeUint64(userBucket, testCoreIntKey, testCoreInt)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	i, err := db.readUint64(userBucket, testCoreIntKey)
	if err != nil {
		t.Fatal("Expected", nil, ", received", err)
	}

	if i != testCoreInt {
		t.Fatal("Expected", testCoreInt, ", received", i)
	}

	db.Close()
}

func testStoreBackup(t *testing.T) {
	fmt.Println(t.Name())

	db := newTestStore(t, testCoreDbPath)

	// Write data to our db
	for _, b := range storeBuckets {
		db.write(b, testCoreKeyName, []byte(testCoreVal1))
		db.writeUint64(b, testCoreIntKey, testCoreInt)
	}

	db.Backup(testCoreBackupDbPath)
	db.Close()

	db2 := newTestStore(t, testCoreBackupDbPath)
	// Read data from backup
	for _, b := range storeBuckets {
		val := db2.read(b, testCoreKeyName)
		if string(val) != testCoreVal1 {
			t.Fatal("Expected", testCoreVal1, ", received", val)
		}

		i, err := db2.readUint64(b, testCoreIntKey)
		if err != nil {
			t.Fatal("Expected", nil, ", received", err)
		}
		if i != testCoreInt {
			t.Fatal("Expected", testCoreInt, ", received", i)
		}
	}

	db2.Close()

	deleteTestStore(t, testCoreDbPath)
	deleteTestStore(t, testCoreBackupDbPath)
}
