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

// KeyValueStorage save key value data
type KeyValueStorage struct {
	data     map[string]map[string]string
	mutex    sync.Mutex
	savePath string
}

func (kv *KeyValueStorage) getBaseDir() string {
	if _, err := os.Stat(kv.savePath); os.IsNotExist(err) {
		os.MkdirAll(kv.savePath, os.ModeDir)
	}
	return kv.savePath
}

func (kv *KeyValueStorage) getLastRecord() (string, error) {
	files, err := ioutil.ReadDir(kv.getBaseDir())
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
	if lastName == "" {
		return "", nil
	}
	return path.Join(kv.getBaseDir(), lastName), nil
}

// Save save data
func (kv *KeyValueStorage) Save() (err error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	kv.mutex.Lock()
	err = enc.Encode(kv.data)
	kv.mutex.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	databaseFile, err := os.Create(path.Join(kv.getBaseDir(), time.Now().Format(dateTimeFilename)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer databaseFile.Close()
	_, err = databaseFile.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	return
}

// Load data
func (kv *KeyValueStorage) Load() (err error) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	lastRecord, err := kv.getLastRecord()
	if err != nil {
		fmt.Println(err)
		return
	}
	if lastRecord == "" {
		return nil
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
	return
}

// AutoInit create new mapping for data, only if needed
func (kv *KeyValueStorage) AutoInit() {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	if len(kv.data) == 0 {
		kv.data = make(map[string]map[string]string)
	}
}

// Set set new value with key
func (kv *KeyValueStorage) Set(k string, v map[string]string) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[k] = v
}

// Get get value with key
func (kv *KeyValueStorage) Get(k string) map[string]string {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	v := kv.data[k]
	return v
}

// List get all keys
func (kv *KeyValueStorage) List() []string {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	keys := make([]string, 0, 200)
	for k := range kv.data {
		keys = append(keys, k)
	}
	return keys
}

// NewKVStorage create or load KeyValueStorage and run persist routine
func NewKVStorage(saveFreq time.Duration, savePathFolderName string) (*KeyValueStorage, chan struct{}, error) {
	savePathFolderName, err := filepath.Abs(savePathFolderName)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	storage := &KeyValueStorage{
		savePath: savePathFolderName,
	}
	err = storage.Load()
	if err != nil {
		fmt.Println(err)
		return storage, nil, err
	}
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
	return storage, quit, nil
}
