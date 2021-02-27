package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const dateTimeFilename = "2006_01_02_15.gob"

func getBaseDir() string {
	path, err := filepath.Abs("./+kvsdata")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}
	return path
}

func getLastRecord() (string, error) {
	files, err := ioutil.ReadDir(getBaseDir())
	if err != nil {
		return "", err
	}
	lastName := ""
	lastTime := time.Unix(1, 0)
	for _, f := range files {
		fTime, err := time.Parse(f.Name(), dateTimeFilename)
		if err != nil {
			if !fTime.After(lastTime) {
				lastTime = fTime
				lastName = f.Name()
			}
		}
	}
	return path.Join(getBaseDir(), lastName), nil
}

// KeyValueStorage save key value data
type KeyValueStorage struct {
	data  map[string]map[string]string
	mutex sync.Mutex
}

// Save save data
func (kv *KeyValueStorage) Save() {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	kv.mutex.Lock()
	err := enc.Encode(kv.data)
	kv.mutex.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	databaseFile, err := os.Create(path.Join(getBaseDir(), time.Now().Format(dateTimeFilename)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer databaseFile.Close()
	_, err = databaseFile.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Load data
func (kv *KeyValueStorage) Load() {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	lastRecord, err := getLastRecord()
	if err != nil {
		fmt.Println(err)
		return
	}
	databaseFile, err := ioutil.ReadFile(lastRecord)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bytes.NewBuffer(databaseFile)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&kv.data)
	if err != nil {
		fmt.Println(err)
	}
}

// AutoInit create new mapping for data, only if needed
func (kv *KeyValueStorage) AutoInit() {
	kv.mutex.Lock()
	if len(kv.data) == 0 {
		kv.data = make(map[string]map[string]string)
	}
	kv.mutex.Unlock()
}

// Set set new value with key
func (kv *KeyValueStorage) Set(k string, v map[string]string) {
	kv.mutex.Lock()
	kv.data[k] = v
	kv.mutex.Unlock()
}

// Get get value with key
func (kv *KeyValueStorage) Get(k string) map[string]string {
	kv.mutex.Lock()
	v := kv.data[k]
	kv.mutex.Unlock()
	return v
}

// List get all keys
func (kv *KeyValueStorage) List() []string {
	kv.mutex.Lock()
	keys := make([]string, 0, 200)
	for k := range kv.data {
		keys = append(keys, k)
	}
	kv.mutex.Unlock()
	return keys
}

// NewKVStorage create or load KeyValueStorage and run persist routine
func NewKVStorage(saveFreq time.Duration) (*KeyValueStorage, chan struct{}) {
	storage := &KeyValueStorage{}
	storage.Load()
	storage.AutoInit()
	persistentRoutine := time.NewTicker(saveFreq)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-persistentRoutine.C:
				fmt.Println("Saving database...")
				go storage.Save()
			case <-quit:
				storage.Save()
				persistentRoutine.Stop()
				return
			}
		}
	}()
	return storage, quit
}
