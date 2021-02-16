package storage

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const passwordSymbols = "abcdefghijklmnopqrstuvwxyz0123456789"
const passwordLength = 12

func generatePassword() string {
	rand.Seed(time.Now().Unix())
	var b strings.Builder
	chars := []rune(passwordSymbols)
	for i := 0; i < passwordLength; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// CreateUser create a new user
func (db *Storage) CreateUser(name string) (*User, error) {
	password := generatePassword()
	tokens := make(map[string]string, 1)
	tokens["general"] = generatePassword()
	user := &User{
		BucketKey:   []byte(name),
		Name:        name,
		BillingPlan: 0,
		Password:    password,
		Tokens:      tokens,
		Visits:      0,
		LastVisit:   time.Now(),
		SystemData: SystemData{
			Created:   time.Now(),
			Updated:   time.Now(),
			IsDeleted: false,
			IsActive:  true,
		},
	}
	savePath := [][]byte{
		[]byte("users"),
	}
	err := db.Set(user, savePath, user.BucketKey)
	return user, err
}

// CreateDomain create a new domain
func (db *Storage) CreateDomain(user *User, hostname string) (*Domain, error) {
	domain := &Domain{
		BucketKey: []byte(hostname),
		Hostname:  hostname,
		SystemData: SystemData{
			Created:   time.Now(),
			Updated:   time.Now(),
			IsDeleted: false,
			IsActive:  true,
		},
	}
	savePath := [][]byte{
		user.BucketKey,
		[]byte("domains"),
	}
	err := db.Set(domain, savePath, domain.BucketKey)
	return domain, err
}

// CreateRedirect create a new redirect
func (db *Storage) CreateRedirect(user *User, domain *Domain, from string, to string) (*Redirect, error) {
	fromURL, err := url.Parse(from)
	if err != nil {
		return nil, err
	}
	toURL, err := url.Parse(to)
	if err != nil {
		return nil, err
	}
	redirect := &Redirect{
		FromURL: *fromURL,
		ToURL:   *toURL,
		SystemData: SystemData{
			Created:   time.Now(),
			Updated:   time.Now(),
			IsDeleted: false,
			IsActive:  true,
		},
	}
	savePath := [][]byte{
		user.BucketKey,
		domain.BucketKey,
		[]byte("redirects"),
	}
	err = db.Set(redirect, savePath, domain.BucketKey)
	return redirect, err
}
