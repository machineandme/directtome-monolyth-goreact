package repository_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/Flaque/filet"
	"github.com/machineandme/directtome-monolyth-goreact/pkg/repository"
	"github.com/stretchr/testify/assert"
)

// TestNew test repository creation
func TestNew(t *testing.T) {
	assert := assert.New(t)
	defer filet.CleanUp(t)
	tf := filet.TmpDir(t, "")
	_, _, err := repository.NewKVStorage(1*time.Second, tf)
	assert.NoError(err)
}

// TestQuit test is quit channel stops persistance routine
func TestQuit(t *testing.T) {
	assert := assert.New(t)
	defer filet.CleanUp(t)
	tf := filet.TmpDir(t, "")
	atStart := runtime.NumGoroutine()
	assert.EqualValues(3, atStart)
	_, quit, _ := repository.NewKVStorage(1*time.Second, tf)
	afterStart := runtime.NumGoroutine()
	assert.EqualValues(4, afterStart)
	close(quit)
	quited := false
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		if runtime.NumGoroutine() == 3 {
			quited = true
			break
		}
	}
	assert.True(quited)
}
