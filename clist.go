package clist

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type IntList struct {
	head   *intNode
	length int64
}

type intNode struct {
	value  int
	next   *intNode
	marked bool
	mu     sync.Mutex
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}

func (n *intNode) loadNext() *intNode {
	return (*intNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *intNode) storeNext(node *intNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(node))
}

func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

func (l *IntList) Insert(value int) bool {
	for {
		// Step 1. Find node a and b
		a := l.head
		b := a.loadNext()
		for b != nil && b.value < value {
			a = b
			b = b.loadNext()
		}
		// Check if the node is exists.
		if b != nil && b.value == value {
			return false
		}
		// Step 2. Lock a
		a.mu.Lock()
		if a.next != b {
			// Step 3. check a.next == b
			a.mu.Unlock()
			continue
		}
		// Step 4. Create new node x
		x := newIntNode(value)
		// Step 5. x.next = b
		x.next = b
		// Step 6. a.next = x
		a.storeNext(x)
		// Step 7. Unlock a
		a.mu.Unlock()
		atomic.AddInt64(&l.length, 1)
		return true
	}
}

func (l *IntList) Delete(value int) bool {
	for {
		// Step 1. Find node a and b
		a := l.head
		b := a.loadNext()
		for b != nil && b.value < value {
			a = b
			b = b.loadNext()
		}

		// Check if b is not exists
		if b == nil {
			return false
		}
		// Step 2. Lock b
		b.mu.Lock()
		if b.value != value || b.marked {
			// Step 3. Check if b has been deleted or another goroutine has delete it
			b.mu.Unlock()
			return false
		}

		// Step 4. Lock a
		a.mu.Lock()
		if a.next != b || a.marked {
			// Step 5. check a.next == b and a is not marked
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}
		// Step 6. mark this node and delete it
		b.marked = true
		a.storeNext(b.next)
		atomic.AddInt64(&l.length, -1)
		// Step 7. unlock a and b
		a.mu.Unlock()
		b.mu.Unlock()
		return true
	}
}

func (l *IntList) Contains(value int) bool {
	x := l.head.loadNext()
	for x != nil && x.value < value {
		x = x.loadNext()
	}
	if x == nil {
		return false
	}
	return x.value == value && !x.marked
}

func (l *IntList) Range(f func(value int) bool) {
	x := l.head.loadNext()
	for x != nil {
		if x.marked {
			x = x.loadNext()
			continue
		}
		if !f(x.value) {
			break
		}
		x = x.loadNext()
	}
}

func (l *IntList) Len() int {
	return int(l.length)
}
