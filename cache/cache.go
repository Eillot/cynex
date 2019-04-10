package cache

import (
	"errors"
	"sync"
)

// 默认容量343
const capability = 7 * 7 * 7

// NewCache 创建缓存，
// 使用LRU算法驱逐缓存中元素
func NewCache(cap ...int) *Cache {
	if len(cap) == 0 {
		return &Cache{
			Cap:      capability,
			Size:     0,
			queue:    &linkedList{head: nil, tail: nil},
			hashData: make(map[string]*node, capability),
			mu:       sync.RWMutex{},
		}
	} else {
		return &Cache{
			Cap:      cap[0],
			Size:     0,
			queue:    &linkedList{head: nil, tail: nil},
			hashData: make(map[string]*node, cap[0]),
			mu:       sync.RWMutex{},
		}
	}
}

type Cache struct {
	Cap  int // 容量
	Size int // 当前存储量

	queue    *linkedList
	hashData map[string]*node
	mu       sync.RWMutex
}

type node struct {
	key  string
	val  interface{}
	prev *node
	next *node
}

type linkedList struct {
	head *node
	tail *node
}

func (l *linkedList) isEmpty() bool {
	if l.head == nil && l.tail == nil {
		return true
	} else {
		return false
	}
}

func (l *linkedList) removeLast() {
	if l.tail != nil {
		l.remove(l.tail)
	}
}

func (l *linkedList) remove(n *node) {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
		return
	}
	if n == l.head {
		n.next.prev = nil
		l.head = n.next
		return
	}
	if n == l.tail {
		n.prev.next = nil
		l.tail = n.prev
		return
	}
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (l *linkedList) addFirst(n *node) {
	if l.head == nil {
		l.head = n
		l.tail = n
		n.prev = nil
		n.next = nil
		return
	}
	n.next = l.head
	l.head.prev = n
	l.head = n
	n.prev = nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if n, ok := c.hashData[key]; ok {
		c.queue.remove(n)
		c.queue.addFirst(n)
		return n.val, nil
	}
	return "", errors.New("not exist")
}

func (c *Cache) Set(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n, ok := c.hashData[key]; ok {
		c.queue.remove(n)
		n.val = val
		c.queue.addFirst(n)
	} else {
		n := &node{key: key, val: val, prev: nil, next: nil}
		c.hashData[key] = n
		c.queue.addFirst(n)
		c.Size += 1
		if c.Size > c.Cap {
			c.Size -= 1
			delete(c.hashData, c.queue.tail.key)
			c.queue.removeLast()
		}
	}
}
