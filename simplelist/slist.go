package simplelist

import (
	"sync"
)

type IntList struct {
	head   *intNode
	length int64
	mu     sync.RWMutex
}

type intNode struct {
	value int
	next  *intNode
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}

func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

func (l *IntList) Insert(value int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.head
	b := a.next
	for b != nil && b.value < value {
		a = b
		b = b.next
	}
	// Check if the node is exist.
	if b != nil && b.value == value {
		return false
	}
	x := newIntNode(value)
	x.next = b
	a.next = x
	l.length++
	return true
}

func (l *IntList) Delete(value int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.head
	b := a.next
	for b != nil && b.value < value {
		a = b
		b = b.next
	}
	// Check if b is not exists
	if b == nil || b.value != value {
		return false
	}
	a.next = b.next
	l.length--
	return true
}

func (l *IntList) Contains(value int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	x := l.head.next
	for x != nil && x.value < value {
		x = x.next
	}
	if x == nil {
		return false
	}
	return x.value == value
}

func (l *IntList) Range(f func(value int) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	x := l.head.next
	for x != nil {
		if !f(x.value) {
			break
		}
		x = x.next
	}
}

func (l *IntList) Len() int {
	return int(l.length)
}
