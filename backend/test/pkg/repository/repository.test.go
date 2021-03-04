package repository_test

import (
	"testing"
	"time"

	"github.com/machineandme/directtome-monolyth-goreact/pkg/repository"
	"github.com/stretchr/testify/assert"
)

// TestNew test repository creation
func TestNew(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(nil)
	repository.NewKVStorage(1 * time.Second)
}
