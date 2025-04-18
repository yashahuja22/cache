package store

import (
	"encoding/json"
	"sync"
)

var (
	once     sync.Once
	instance DataStore
)

type DataStore interface {
	JSONGet(key string) (json.RawMessage, bool)
	JSONSet(key string, jsonData json.RawMessage)
	JSONDel(key string) bool
}

type storageManager struct {
	mu   sync.RWMutex
	data map[string]json.RawMessage
}

func (sm *storageManager) JSONGet(key string) (json.RawMessage, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	d, ok := sm.data[key]

	return d, ok

}
func (sm *storageManager) JSONSet(key string, jsonData json.RawMessage) {
	sm.data[key] = jsonData
}
func (sm *storageManager) JSONDel(key string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, exists := sm.data[key]; exists {
		delete(sm.data, key)
		return true
	}
	return false
}

func NewStorageManager() DataStore {
	// singleton design pattern
	once.Do(func() {
		instance = &storageManager{
			data: make(map[string]json.RawMessage),
		}
	})
	return instance
}
