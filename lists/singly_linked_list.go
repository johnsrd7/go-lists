package listadts

import (
	"sync"

	adts "github.com/johnsrd7/go-adts"
)

type listNode struct {
	elt  adts.ContainerElement
	next *listNode
}

// SinglyLinkedList is a simple singly linked list object that
// can be made threadsafe.
type SinglyLinkedList struct {
	head       *listNode
	tail       *listNode
	len        int
	lock       *sync.Mutex
	threadSafe bool
}

// makeListNode creates a new listNode object.
func makeListNode(elt adts.ContainerElement) *listNode {
	return &listNode{elt, nil}
}

// MakeSinglyLinkedList creates a non-threadsafe SinglyLinkedList
func MakeSinglyLinkedList() *SinglyLinkedList {
	return &SinglyLinkedList{nil, nil, 0, &sync.Mutex{}, false}
}

// MakeSinglyLinkedListThreadsafe creates a new SinglyLinkedList that is threadsafe.
func MakeSinglyLinkedListThreadsafe() *SinglyLinkedList {
	return &SinglyLinkedList{nil, nil, 0, &sync.Mutex{}, true}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the list.
func (l *SinglyLinkedList) Len() int {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.len
	}
	return l.len
}

// IsEmpty returns if the list is empty or not.
func (l *SinglyLinkedList) IsEmpty() bool {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.len == 0
	}

	return l.len == 0
}

// Clear removes all elements from the list.
func (l *SinglyLinkedList) Clear() {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		l.head = nil
		l.tail = nil
		l.len = 0
	}

	l.head = nil
	l.tail = nil
	l.len = 0
}

// Contains returns true if the given item is in the list.
func (l *SinglyLinkedList) Contains(item adts.ContainerElement) bool {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.containsHelper(item)
	}

	return l.containsHelper(item)
}

// containsHelper returns whether the given element is in the list.
func (l *SinglyLinkedList) containsHelper(item adts.ContainerElement) bool {
	if l.head == nil {
		return false
	}

	for tmp := l.head; tmp != nil; tmp = tmp.next {
		if tmp.elt.Equals(item) {
			return true
		}
	}

	return false
}

// Add returns true if the given element was appended to the end of the list.
func (l *SinglyLinkedList) Add(item adts.ContainerElement) bool {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.addHelper(item)
	}

	return l.addHelper(item)
}

// addHelper add the given element to the end of the linked list.
func (l *SinglyLinkedList) addHelper(item adts.ContainerElement) bool {
	newNode := makeListNode(item)

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = l.tail.next
	}

	// Don't forget to update the length.
	l.len++
	return true
}

// Remove removes the element at the given index and
// returns whether the removal was successful.
func (l *SinglyLinkedList) Remove(item adts.ContainerElement) bool {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.removeHelper(item)
	}

	return l.removeHelper(item)
}

// removeHelper searches the list for the given element and then just
// sets the next links properly to remove the element from the list.
func (l *SinglyLinkedList) removeHelper(item adts.ContainerElement) bool {
	if l.head == nil {
		return false
	}

	// Check if the head is the item to be removed
	if l.head.elt.Equals(item) {
		l.head = l.head.next
		l.len--
		return true
	}

	for tmp := l.head; tmp.next != nil; tmp = tmp.next {
		if tmp.next.elt.Equals(item) {
			tmp.next = tmp.next.next
			// Don't forget to update the length.
			l.len--
			return true
		}
	}

	return false
}

// -------------------------------------------------------
// List Methods
// -------------------------------------------------------

// Get returns the element at the given index.
func (l *SinglyLinkedList) Get(idx int) adts.ContainerElement {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.getHelper(idx)
	}

	return l.getHelper(idx)
}

// getHelper searches the list for the element and returns it if the
// index is within the list.
func (l *SinglyLinkedList) getHelper(idx int) adts.ContainerElement {
	curIdx := 0
	for tmp := l.head; tmp != nil; tmp = tmp.next {
		if curIdx == idx {
			return tmp.elt
		}

		curIdx++
	}

	panic("index out of range")
}

// Set changes the value at the given index to the given new value
// and returns the old value that was at the given index.
func (l *SinglyLinkedList) Set(idx int, newVal adts.ContainerElement) adts.ContainerElement {
	if l.threadSafe {
		l.lock.Lock()
		defer l.lock.Unlock()
		return l.setHelper(idx, newVal)
	}

	return l.setHelper(idx, newVal)
}

// setHelper searches the list for the element and returns it if the
// index is within the list.
func (l *SinglyLinkedList) setHelper(idx int, newVal adts.ContainerElement) adts.ContainerElement {
	curIdx := 0
	for tmp := l.head; tmp != nil; tmp = tmp.next {
		if curIdx == idx {
			oldVal := tmp.elt
			tmp.elt = newVal
			return oldVal
		}

		curIdx++
	}

	panic("index out of range")
}
