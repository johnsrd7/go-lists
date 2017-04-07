package stackadts

import (
	adts "github.com/johnsrd7/go-adts"
)

// SliceStack is a simple type that implements the Stack interface (both threadsafe and not).
type SliceStack struct {
	backer *adts.SliceContainer
}

// MakeSliceStack creates a non-threadsafe SliceStack
func MakeSliceStack() *SliceStack {
	return &SliceStack{adts.MakeSliceContainer()}
}

// MakeSliceStackThreadSafe creates a threadsafe SliceStack
func MakeSliceStackThreadSafe() *SliceStack {
	return &SliceStack{adts.MakeSliceContainerThreadSafe()}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the stack.
func (ss *SliceStack) Len() int {
	return ss.backer.Len()
}

// IsEmpty returns if the stack is empty or not.
func (ss *SliceStack) IsEmpty() bool {
	return ss.backer.Len() == 0
}

// Clear removes all elements from the stack.
func (ss *SliceStack) Clear() {
	ss.backer.Clear()
}

// Contains returns true if the given item is in the stack.
func (ss *SliceStack) Contains(item adts.ContainerElement) bool {
	return ss.backer.Contains(item)
}

// Add returns true if the given element was added to the top of the stack.
func (ss *SliceStack) Add(item adts.ContainerElement) bool {
	return ss.backer.Add(item)
}

// Remove is a non-valid function for the Stack interface. It is only provided here
// as a means to satisfy the Container interface. See Pop function for Stack to remove
// elements from the stack.
func (ss *SliceStack) Remove(item adts.ContainerElement) bool {
	return ss.backer.Remove(item)
}

// -------------------------------------------------------
// Stack Methods
// -------------------------------------------------------

// Push pushes the given element onto the top of the stack.
func (ss *SliceStack) Push(item adts.ContainerElement) {
	if !ss.Add(item) {
		panic("Unable to push the given item to the top of the stack.")
	}
}

// Pop removes the top element from the stack and returns the element.
func (ss *SliceStack) Pop() adts.ContainerElement {
	// We want to reuse the Remove method from the SliceContainer class.
	// The problem is that we need to get the last element in a threadsafe way
	// (if needed) and then call remove. However, if we lock and then call
	// Remove, that will also lock, which causes a deadlock. So we need to
	// do our own locking here and then remove the element and then unlock.

	var lastElt adts.ContainerElement
	if ss.backer.ThreadSafe {
		ss.backer.Lock.Lock()
		lastElt = ss.backer.Backer[len(ss.backer.Backer)-1]
		ss.backer.RemoveAtIndex(len(ss.backer.Backer) - 1)
		ss.backer.Lock.Unlock()
	} else {
		lastElt = ss.backer.Backer[len(ss.backer.Backer)-1]
		ss.backer.RemoveAtIndex(len(ss.backer.Backer) - 1)
	}

	return lastElt
}
