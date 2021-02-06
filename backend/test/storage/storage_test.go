package storage_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/machineandme/directtome-monolyth-goreact/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var tempDir string

func init() {
	rand.Seed(time.Now().UnixNano())
	tempDir = "./data/" + strconv.Itoa(rand.Int())
}

func getDatabasePath(t *testing.T) string {
	tempSubDir := filepath.FromSlash(tempDir + "/" + strconv.Itoa(rand.Int()))
	os.MkdirAll(tempSubDir, 0777)
	databasePath := tempSubDir + "/bolt.db"
	databasePath = filepath.FromSlash(databasePath)
	os.Create(databasePath)
	return databasePath
}

func TestOpenStorage(t *testing.T) {
	databasePath := getDatabasePath(t)
	fileInfo, err := os.Stat(databasePath)
	assert.NoError(t, err)
	assert.EqualValues(t, fileInfo.Size(), 0)
	repo := &storage.Storage{}
	err = repo.OpenDatabase(databasePath, "v0.0.0test")
	assert.NoError(t, err)
	fileInfo, err = os.Stat(databasePath)
	assert.NoError(t, err)
	assert.NotEqualValues(t, fileInfo.Size(), 0)
	repo.CloseDatabase()
}

func TestCloseStorage(t *testing.T) {
	databasePath := getDatabasePath(t)
	repo := &storage.Storage{}
	repo.OpenDatabase(databasePath, "v0.0.0test")
	isOpen := reflect.ValueOf(*repo).FieldByName("isOpen").Bool()
	assert.EqualValues(t, isOpen, true)
	err := repo.CloseDatabase()
	assert.NoError(t, err)
	isOpen = reflect.ValueOf(*repo).FieldByName("isOpen").Bool()
	assert.EqualValues(t, isOpen, false)
}

func TestStoreSomething(t *testing.T) {
	databasePath := getDatabasePath(t)
	repo := &storage.Storage{}
	repo.OpenDatabase(databasePath, "v0.0.0test")
	bucketKeys := [][]byte{
		[]byte("A"),
		[]byte("1"),
	}
	key := []byte("data")
	type dollStruct struct{ ValueAlpha, ValueBetta int }
	settedDoll := dollStruct{
		ValueAlpha: 1,
		ValueBetta: 2,
	}
	fileInfo, err := os.Stat(databasePath)
	assert.NoError(t, err)
	delta := fileInfo.Size()
	err = repo.Set(&settedDoll, bucketKeys, key)
	assert.NoError(t, err)
	fileInfo, err = os.Stat(databasePath)
	assert.NoError(t, err)
	assert.NotEqual(t, fileInfo.Size(), delta)
	readDone := make(chan bool)
	gettedDoll := dollStruct{}
	err = repo.Get(&gettedDoll, bucketKeys, key, readDone)
	assert.NoError(t, err)
	<-readDone
	assert.Equal(t, gettedDoll, settedDoll)
}
