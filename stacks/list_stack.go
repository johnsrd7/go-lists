package stackadts

import (
	"container/list"
	"sync"

	adts "github.com/johnsrd7/go-adts"
)

// ListStack is a simple type that implements the Stack interface (both threadsafe and not).
type ListStack struct {
	backer     *list.List
	lock       *sync.Mutex
	threadSafe bool
}

// MakeListStack creates a non-threadsafe ListStack
func MakeListStack() *ListStack {
	return &ListStack{list.New(), &sync.Mutex{}, false}
}

// MakeListStackThreadSafe creates a threadsafe ListStack
func MakeListStackThreadSafe() *ListStack {
	return &ListStack{list.New(), &sync.Mutex{}, true}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the stack.
func (ls *ListStack) Len() int {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		return ls.backer.Len()
	}

	return ls.backer.Len()
}

// IsEmpty returns if the stack is empty or not.
func (ls *ListStack) IsEmpty() bool {
	return ls.backer.Len() == 0
}

// Clear removes all elements from the stack.
func (ls *ListStack) Clear() {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		ls.backer.Init()
		return
	}

	ls.backer.Init()
}

// Contains returns true if the given item is in the stack.
func (ls *ListStack) Contains(item adts.ContainerElement) bool {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		return ls.containsHelper(item)
	}

	return ls.containsHelper(item)
}

// containsHelper checks the list for the given element (in a non-threadsafe way).
func (ls *ListStack) containsHelper(item adts.ContainerElement) bool {
	for tmp := ls.backer.Front(); tmp != nil; tmp = tmp.Next() {
		if v, ok := tmp.Value.(adts.ContainerElement); ok {
			if v.Equals(item) {
				return true
			}
		}
	}

	return false
}

// Add returns true if the given element was added to the top of the stack.
func (ls *ListStack) Add(item adts.ContainerElement) bool {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		return ls.backer.PushFront(item) != nil
	}

	return ls.backer.PushFront(item) != nil
}

// Remove returns true if the given element was removed.
func (ls *ListStack) Remove(item adts.ContainerElement) bool {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		return ls.removeHelper(item)
	}

	return ls.removeHelper(item)
}

func (ls *ListStack) removeHelper(item adts.ContainerElement) bool {
	for tmp := ls.backer.Front(); tmp != nil; tmp = tmp.Next() {
		if v, ok := tmp.Value.(adts.ContainerElement); ok {
			if v.Equals(item) {
				removed, ok := ls.backer.Remove(tmp).(adts.ContainerElement)
				return ok && removed.Equals(item)
			}
		}
	}

	return false
}

// -------------------------------------------------------
// Stack Methods
// -------------------------------------------------------

// Push pushes the given element onto the top of the stack.
func (ls *ListStack) Push(item adts.ContainerElement) bool {
	return ls.Add(item)
}

// Pop removes the top element from the stack and returns the element.
func (ls *ListStack) Pop() (adts.ContainerElement, bool) {
	if ls.threadSafe {
		ls.lock.Lock()
		defer ls.lock.Unlock()
		return ls.popHelper()
	}

	return ls.popHelper()
}

// popHelper removes the element from the front of the list and returns the element.
func (ls *ListStack) popHelper() (adts.ContainerElement, bool) {
	// Have to use list.List's Len function since we might already
	// be inside a lock and we don't want to deadlock.
	if ls.backer.Len() == 0 {
		return adts.EmptyContainerElement{}, false
	}

	if lastElt, ok := ls.backer.Front().Value.(adts.ContainerElement); ok {
		ls.backer.Remove(ls.backer.Front())
		return lastElt, true
	}

	return adts.EmptyContainerElement{}, false
}
