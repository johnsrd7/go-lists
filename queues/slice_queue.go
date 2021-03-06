package queueadts

import (
	adts "github.com/johnsrd7/go-adts"
)

// SliceQueue is a simple type that implements the Queue interface (both threadsafe and not).
type SliceQueue struct {
	backer *adts.SliceContainer
}

// MakeSliceQueue creates a non-threadsafe SliceQueue
func MakeSliceQueue() *SliceQueue {
	return &SliceQueue{adts.MakeSliceContainer()}
}

// MakeSliceQueueThreadSafe creates a threadsafe SliceQueue
func MakeSliceQueueThreadSafe() *SliceQueue {
	return &SliceQueue{adts.MakeSliceContainerThreadSafe()}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the queue.
func (sq *SliceQueue) Len() int {
	return sq.backer.Len()
}

// IsEmpty returns if the queue is empty or not.
func (sq *SliceQueue) IsEmpty() bool {
	return sq.backer.Len() == 0
}

// Clear removes all elements from the queue.
func (sq *SliceQueue) Clear() {
	sq.backer.Clear()
}

// Contains returns true if the given item is in the queue.
func (sq *SliceQueue) Contains(item adts.ContainerElement) bool {
	return sq.backer.Contains(item)
}

// Add returns true if the given element was added to the end of the queue.
func (sq *SliceQueue) Add(item adts.ContainerElement) bool {
	return sq.backer.Add(item)
}

// Remove returns true if the given element was removed.
func (sq *SliceQueue) Remove(item adts.ContainerElement) bool {
	return sq.backer.Remove(item)
}

// -------------------------------------------------------
// Stack Methods
// -------------------------------------------------------

// Enqueue pushes the given element onto the back of the queue.
func (sq *SliceQueue) Enqueue(item adts.ContainerElement) bool {
	return sq.Add(item)
}

// Dequeue removes the element from the head of the queue and returns the element.
func (sq *SliceQueue) Dequeue() (adts.ContainerElement, bool) {
	// We want to reuse the Remove method from the SliceContainer class.
	// The problem is that we need to get the last element in a threadsafe way
	// (if needed) and then call remove. However, if we lock and then call
	// Remove, that will also lock, which causes a deadlock. So we need to
	// do our own locking here and then remove the element and then unlock.
	if sq.Len() == 0 {
		return adts.EmptyContainerElement{}, false
	}

	var firstElt adts.ContainerElement
	if sq.backer.ThreadSafe {
		sq.backer.Lock.Lock()
		defer sq.backer.Lock.Unlock()
		firstElt = sq.backer.Backer[0]
		if !sq.backer.RemoveAtIndex(0) {
			return adts.EmptyContainerElement{}, false
		}
	} else {
		firstElt = sq.backer.Backer[0]
		if !sq.backer.RemoveAtIndex(0) {
			return adts.EmptyContainerElement{}, false
		}
	}

	return firstElt, true
}
