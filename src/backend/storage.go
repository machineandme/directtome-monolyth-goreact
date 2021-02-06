package main

import (
	"bytes"
	"encoding/gob"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
)

var db bolt.DB

// Database name `Users`. Path key style `{Email}`
type User struct {
	Active        bool
	Name          string
	Password      string
	Token         string
	Email         string
	Phone         string
	Company       string
	Country       string
	City          string
	LastUA        string
	Visits        uint64
	LastVisit     time.Time
	TimeZone      time.Location
	RootDomainKey string
}

// Database name `Domains`. Path key style `{Hostname}`
type Domain struct {
	Active                 bool
	Hostname               string
	NotFoundUrl            string
	FeatureGeoIP           bool
	FeatureAutoHTTPS       bool
	FeatureCaptureUA       bool
	FeatureCaptureReferrer bool
	StatKey                string
	NextDomainKey          string
	RootRedirectKey        string
}

// Database name `Redirects`. Path key style `{FromHost} {FromPath}`
type Redirect struct {
	Active          bool
	FromHost        string
	FromPath        string
	To              url.URL
	Canonical       url.URL
	StatusCode      uint8
	CacheControl    string
	ExpireClicks    uint64
	ExpireTime      time.Time
	StatKey         string
	NextRedirectKey string
}

// Database name `Statistics`. Path key is random value base64 encoded
type Statistic struct {
	Day               time.Time
	Clicks            uint64
	ExpiredClicks     uint32
	NotFoundClicks    uint32
	SocialShareClicks uint64
	ExtendedStats     map[string]uint64
}

func initDatabase() {
	db, err := bolt.Open("/var/application_backend/main.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Domains"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Redirects"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Statistics"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func stopDatabase() {
	db.Close()
}

func getValue(database string, key string, value *interface{}) {
	err := db.View(func(tx *bolt.Tx) error {
		valueBuffer := tx.Bucket([]byte(database)).Get([]byte(key))
		dec := gob.NewDecoder(&valueBuffer)
		err := dec.Decode(value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func setValue(database string, key string, value interface{}) {
	var valueBuffer bytes.Buffer
	enc := gob.NewEncoder(&valueBuffer)
	err := enc.Encode(value)
	if err != nil {
		panic(err)
	}
	valueBytes := valueBuffer.Bytes()
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(database)).Put([]byte(key), valueBytes)
	})
	if err != nil {
		panic(err)
	}
	return
}

func createAdmin() User {
	admin := User{}
	setValue("Users", "admin@local", User{
		Active:    true,
		Name:      "Admin",
		Password:  "Admin",
		Token:     "Admin",
		Email:     "admin@local",
		LastVisit: time.Now(),
		TimeZone:  *time.FixedZone("UTC+3", 3*60*60),
	})
	getValue("Users", "admin@local", &admin)
	return admin
}
