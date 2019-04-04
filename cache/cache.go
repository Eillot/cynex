package cache

import "errors"

// 默认容量
const capability = 343

// 使用LRU算法存储缓存
func NewCache(cap ...int) *Cache {
	if len(cap) == 0 {
		return &Cache{
			Cap:     capability,
			Size:    0,
			Queue:   &LinkedList{Head: nil, Tail: nil},
			HashMap: make(map[string]*Node, capability),
		}
	} else {
		return &Cache{
			Cap:     cap[0],
			Size:    0,
			Queue:   &LinkedList{Head: nil, Tail: nil},
			HashMap: make(map[string]*Node, cap[0]),
		}
	}
}

type Cache struct {
	Cap     int
	Size    int
	Queue   *LinkedList
	HashMap map[string]*Node
}

type Node struct {
	Key  string
	Val  string
	Prev *Node
	Next *Node
}

type LinkedList struct {
	Head *Node
	Tail *Node
}

func (l *LinkedList) IsEmpty() bool {
	if l.Head == nil && l.Tail == nil {
		return true
	} else {
		return false
	}
}

func (l *LinkedList) RemoveLast() {
	if l.Tail != nil {
		l.Remove(l.Tail)
	}
}

func (l *LinkedList) Remove(n *Node) {
	if l.Head == l.Tail {
		l.Head = nil
		l.Tail = nil
		return
	}
	if n == l.Head {
		n.Next.Prev = nil
		l.Head = n.Next
		return
	}
	if n == l.Tail {
		n.Prev.Next = nil
		l.Tail = n.Prev
		return
	}
	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
}

func (l *LinkedList) AddFirst(n *Node) {
	if l.Head == nil {
		l.Head = n
		l.Tail = n
		n.Prev = nil
		n.Next = nil
		return
	}
	n.Next = l.Head
	l.Head.Prev = n
	l.Head = n
	n.Prev = nil
}

func (c *Cache) Get(key string) (string, error) {
	if node, ok := c.HashMap[key]; ok {
		c.Queue.Remove(node)
		c.Queue.AddFirst(node)
		return node.Val, nil
	}
	return "", errors.New("not exist")
}

func (c *Cache) Set(key, val string) {
	if node, ok := c.HashMap[key]; ok {
		c.Queue.Remove(node)
		node.Val = val
		c.Queue.AddFirst(node)
	} else {
		n := &Node{Key: key, Val: val, Prev: nil, Next: nil}
		c.HashMap[key] = n
		c.Queue.AddFirst(n)
		c.Size += 1
		if c.Size > c.Cap {
			c.Size -= 1
			delete(c.HashMap, c.Queue.Tail.Key)
			c.Queue.RemoveLast()
		}
	}
}
