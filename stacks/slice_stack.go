package stackadts

import (
	"sync"

	adts "github.com/johnsrd7/go-adts"
)

// SliceStack is a simple type that implements the Stack interface (both threadsafe and not).
type SliceStack struct {
	backer       []adts.ContainerElement
	lock         *sync.Mutex
	threadSafe   bool
	shrinkFactor float32
}

// MakeSliceStack creates a non-threadsafe SliceStack
func MakeSliceStack() *SliceStack {
	return &SliceStack{[]adts.ContainerElement{}, &sync.Mutex{}, false, 0.25}
}

// MakeSliceStackThreadSafe creates a threadsafe SliceStack
func MakeSliceStackThreadSafe() *SliceStack {
	return &SliceStack{[]adts.ContainerElement{}, &sync.Mutex{}, true, 0.25}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the stack.
func (ss *SliceStack) Len() int {
	if ss.threadSafe {
		ss.lock.Lock()
		defer ss.lock.Unlock()
		return len(ss.backer)
	}
	return len(ss.backer)
}

// IsEmpty returns if the stack is empty or not.
func (ss *SliceStack) IsEmpty() bool {
	return ss.Len() == 0
}

// Clear removes all elements from the stack.
func (ss *SliceStack) Clear() {
	if ss.threadSafe {
		ss.lock.Lock()
		defer ss.lock.Unlock()
		ss.backer = []adts.ContainerElement{}
		return
	}

	ss.backer = []adts.ContainerElement{}
}

// Contains returns true if the given item is in the stack.
func (ss *SliceStack) Contains(item adts.ContainerElement) bool {
	if ss.threadSafe {
		ss.lock.Lock()
		defer ss.lock.Unlock()
		return ss.containsHelper(item)
	}

	return ss.containsHelper(item)
}

func (ss *SliceStack) containsHelper(item adts.ContainerElement) bool {
	for _, i := range ss.backer {
		if i.Equals(item) {
			return true
		}
	}

	return false
}

// Add returns true if the given element was added to the top of the stack.
func (ss *SliceStack) Add(item adts.ContainerElement) bool {
	ss.Push(item)
	return true
}

// Remove is a non-valid function for the Stack interface. It is only provided here
// as a means to satisfy the Container interface. See Pop function for Stack to remove
// elements from the stack.
func (ss *SliceStack) Remove(item adts.ContainerElement) bool {
	return false
}

// -------------------------------------------------------
// List Methods
// -------------------------------------------------------

// Push pushes the given element onto the top of the stack.
func (ss *SliceStack) Push(item adts.ContainerElement) {
	if ss.threadSafe {
		ss.lock.Lock()
		defer ss.lock.Unlock()
		ss.backer = append(ss.backer, item)
		return
	}

	ss.backer = append(ss.backer, item)
}

// Pop removes the top element from the stack and returns the element.
func (ss *SliceStack) Pop() adts.ContainerElement {
	if ss.threadSafe {
		ss.lock.Lock()
		defer ss.lock.Unlock()
		return ss.popHelper()
	}

	return ss.popHelper()
}

func (ss *SliceStack) popHelper() adts.ContainerElement {
	lastElt := ss.backer[len(ss.backer)-1]
	ss.backer = ss.backer[:len(ss.backer)-1]

	// We should check to see if we need to resize the slice. We don't
	// want it to be the case that we added a ton of items then removed
	// a bunch and now we are still holding onto the large backing array
	// for the slice.
	emptyFactor := float32(len(ss.backer)) / float32(cap(ss.backer))
	if emptyFactor <= ss.shrinkFactor {
		// Shrink the slice's capacity by 1/2
		newCap := cap(ss.backer) / 2
		newBacker := make([]adts.ContainerElement, len(ss.backer), newCap)
		copy(newBacker, ss.backer)
		ss.backer = newBacker
	}

	return lastElt
}
