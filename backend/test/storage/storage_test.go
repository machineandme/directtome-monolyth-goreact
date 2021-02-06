package storage_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/machineandme/directtome-monolyth-goreact/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func getDatabasePath(t *testing.T) string {
	databasePath := "." + "/bolt.db"
	return databasePath
}

func TestOpenStorage(t *testing.T) {
	databasePath := getDatabasePath(t)
	_, err := os.Stat(databasePath)
	assert.Error(t, err)
	repo := &storage.Storage{}
	err = repo.OpenDatabase(databasePath, "v0.0.0test")
	assert.NoError(t, err)
	fileInfo, err := os.Stat(databasePath)
	assert.NoError(t, err)
	assert.NotEqualValues(t, fileInfo.Size(), 0)
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
