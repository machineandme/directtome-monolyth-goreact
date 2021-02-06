package storage

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/boltdb/bolt"
)

// Storage is structure that helps to communicate with database
type Storage struct {
	isOpen        bool
	DatabasePath  string
	BucketVersion []byte
	writeMutex    sync.Mutex
	readWaitGroup sync.WaitGroup
	boltDB        *bolt.DB
}

func (db *Storage) open() error {
	db.isOpen = true
	boltDB, err := bolt.Open(db.DatabasePath, 0777, nil)
	if err != nil {
		return err
	}
	db.boltDB = boltDB
	return nil
}

func (db *Storage) close() error {
	db.isOpen = false
	err := db.boltDB.Close()
	return err
}

// OpenDatabase lock for operation and open file with bold database
func (db *Storage) OpenDatabase(path string, bucketVersion string) error {
	db.writeMutex.Lock()
	if db.isOpen {
		err := db.close()
		if err != nil {
			return err
		}
	}
	db.DatabasePath = path
	db.BucketVersion = []byte(bucketVersion)
	err := db.open()
	if err != nil {
		return err
	}
	err = db.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(db.BucketVersion)
		return err
	})
	db.writeMutex.Unlock()
	return err
}

// CloseDatabase wait for operation complete and close file with bolt database
func (db *Storage) CloseDatabase() error {
	if db.isOpen {
		db.writeMutex.Lock()
		db.readWaitGroup.Wait()
		err := db.close()
		db.writeMutex.Unlock()
		return err
	}
	return nil
}

// Reopen wait for operation complete and reopen file with bolt database
func (db *Storage) Reopen() error {
	err := db.CloseDatabase()
	if err != nil {
		return err
	}
	err = db.Connect()
	return err
}

// Connect lock for operation and open file with bold database if database is closed
func (db *Storage) Connect() error {
	var err error = nil
	if !db.isOpen {
		db.writeMutex.Lock()
		db.readWaitGroup.Wait()
		err = db.open()
		db.writeMutex.Unlock()
	}
	db.boltDB.Sync()
	return err
}

// Get data to model by key
func (db *Storage) Get(dataModel interface{}, bucketKeys [][]byte, key []byte, notify chan bool) error {
	err := db.Connect()
	if err != nil {
		return err
	}
	db.writeMutex.Lock()
	db.readWaitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := db.boltDB.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(db.BucketVersion)
			for _, bucketKey := range bucketKeys {
				bucket = bucket.Bucket(bucketKey)
			}
			rawData := bucket.Get(key)
			rawDataBuff := bytes.NewBuffer(rawData)
			dec := gob.NewDecoder(rawDataBuff)
			return dec.Decode(dataModel)
		})
		notify <- err != nil
		if err != nil {
			panic(err)
		}
	}(&db.readWaitGroup)
	db.writeMutex.Unlock()
	return nil
}

// Set data to model by key
func (db *Storage) Set(dataModel interface{}, bucketKeys [][]byte, key []byte) error {
	err := db.Connect()
	if err != nil {
		return err
	}
	db.writeMutex.Lock()
	err = db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.BucketVersion)
		for _, bucketKey := range bucketKeys {
			var err error
			bucket, err = bucket.CreateBucketIfNotExists(bucketKey)
			if err != nil {
				return err
			}
		}
		var rawDataBuff bytes.Buffer
		enc := gob.NewEncoder(&rawDataBuff)
		err := enc.Encode(dataModel)
		if err != nil {
			return err
		}
		err = bucket.Put(key, rawDataBuff.Bytes())
		return err
	})
	db.writeMutex.Unlock()
	return err
}
