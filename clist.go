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
	marked uint32
	next   *intNode
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

func (n *intNode) setMarked() {
	atomic.StoreUint32(&n.marked, 1)
}

func (n *intNode) isMarked() bool {
	return atomic.LoadUint32(&n.marked) == 1
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
			if b.isMarked() {
				// The node has been marked, insert it in next loop.
				continue
			}
			return false
		}
		// Step 2. Lock a
		a.mu.Lock()
		if a.next != b { // check a.next == b
			// Version 1.
			a.mu.Unlock()
			continue
			// Version 2.
			// if a.next == nil || a.next.value < value {
			// 	a.mu.Unlock()
			// 	continue
			// }
			// if a.next.value == value {
			// 	a.mu.Unlock()
			// 	return false
			// } else if a.next.value > value {
			// 	b = a.next
			// }
		}
		// Step 3. Create new node x
		x := newIntNode(value)
		// Step 4. x.next = b; a.next = x
		x.next = b
		a.storeNext(x)
		// Step 5. Unlock a
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

		// Check if b is not exist.
		if b == nil || b.value != value {
			return false
		}

		// Step 2. Lock b
		b.mu.Lock()
		if b.isMarked() {
			// Check if b has been deleted or another goroutine has delete it.
			b.mu.Unlock()
			return false
		}

		// Step 3. Lock a
		a.mu.Lock()
		if a.next != b || a.isMarked() {
			// Check a.next == b and a is not marked
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}
		// Step 4. Mark this node and delete it
		b.setMarked()
		a.storeNext(b.next)
		atomic.AddInt64(&l.length, -1)
		// Step 5. Unlock a and b
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
	return x.value == value && !x.isMarked()
}

func (l *IntList) Range(f func(value int) bool) {
	x := l.head.loadNext()
	for x != nil {
		if x.isMarked() {
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
	return int(atomic.LoadInt64(&l.length))
}
