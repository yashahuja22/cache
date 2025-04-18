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

type entry struct {
	key   string
	value json.RawMessage
}

type node struct {
	entry *entry
	prev  *node
	next  *node
}

type storageManager struct {
	capacity int
	items    map[string]*node
	head     *node
	tail     *node
	mu       sync.RWMutex
}

func (sm *storageManager) JSONGet(key string) (json.RawMessage, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	d, ok := sm.items[key]
	if ok {
		sm.moveToFront(d)
		return d.entry.value, ok
	}

	return nil, ok

}
func (sm *storageManager) JSONSet(key string, value json.RawMessage) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if node, exists := sm.items[key]; exists {
		node.entry.value = value
		sm.moveToFront(node)
		return
	}

	if len(sm.items) >= sm.capacity {
		// Remove LRU
		delete(sm.items, sm.tail.entry.key)
		sm.removeNode(sm.tail)
	}

	newNode := &node{
		entry: &entry{
			key:   key,
			value: value,
		},
	}

	sm.addToFront(newNode)
	sm.items[key] = newNode
}
func (sm *storageManager) JSONDel(key string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if node, exists := sm.items[key]; exists {
		sm.removeNode(node)
		delete(sm.items, key)
		return true
	}
	return false
}

func (c *storageManager) moveToFront(node *node) {
	if node == c.head {
		return
	}
	c.removeNode(node)
	c.addToFront(node)
}

func (c *storageManager) removeNode(node *node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

func (c *storageManager) addToFront(node *node) {
	node.next = c.head
	node.prev = nil

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func NewStorageManager(capacity int) DataStore {
	// singleton design pattern
	once.Do(func() {
		instance = &storageManager{
			capacity: capacity,
			items:    make(map[string]*node),
		}
	})
	return instance
}
