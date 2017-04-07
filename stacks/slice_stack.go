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
	ss.Push(item)
	return true
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
	if ss.backer.ThreadSafe {
		ss.backer.Lock.Lock()
		defer ss.backer.Lock.Unlock()
		ss.backer.Backer = append(ss.backer.Backer, item)
		return
	}

	ss.backer.Backer = append(ss.backer.Backer, item)
}

// Pop removes the top element from the stack and returns the element.
func (ss *SliceStack) Pop() adts.ContainerElement {
	if ss.backer.ThreadSafe {
		ss.backer.Lock.Lock()
		defer ss.backer.Lock.Unlock()
		return ss.popHelper()
	}

	return ss.popHelper()
}

func (ss *SliceStack) popHelper() adts.ContainerElement {
	lastElt := ss.backer.Backer[len(ss.backer.Backer)-1]
	ss.backer.Backer = ss.backer.Backer[:len(ss.backer.Backer)-1]

	// We should check to see if we need to resize the slice. We don't
	// want it to be the case that we added a ton of items then removed
	// a bunch and now we are still holding onto the large backing array
	// for the slice.
	emptyFactor := float32(len(ss.backer.Backer)) / float32(cap(ss.backer.Backer))
	if emptyFactor <= ss.backer.ShrinkFactor {
		// Shrink the slice's capacity by 1/2
		newCap := cap(ss.backer.Backer) / 2
		newBacker := make([]adts.ContainerElement, len(ss.backer.Backer), newCap)
		copy(newBacker, ss.backer.Backer)
		ss.backer.Backer = newBacker
	}

	return lastElt
}
