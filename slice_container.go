package adts

import (
	"sync"
)

// SliceContainer is a simple type that implements the Container interface.
type SliceContainer struct {
	Backer       []ContainerElement
	Lock         *sync.Mutex
	ThreadSafe   bool
	ShrinkFactor float32
}

// MakeSliceContainer creates a new non-threadsafe SliceContainer.
func MakeSliceContainer() *SliceContainer {
	return &SliceContainer{[]ContainerElement{}, &sync.Mutex{}, false, 0.25}
}

// MakeSliceContainerThreadSafe creates a new threadsafe SliceContainer.
func MakeSliceContainerThreadSafe() *SliceContainer {
	return &SliceContainer{[]ContainerElement{}, &sync.Mutex{}, true, 0.25}
}

// Len returns the number of elements in the container.
func (sc *SliceContainer) Len() int {
	if sc.ThreadSafe {
		sc.Lock.Lock()
		defer sc.Lock.Unlock()
		return len(sc.Backer)
	}
	return len(sc.Backer)
}

// IsEmpty returns if the container is empty or not.
func (sc *SliceContainer) IsEmpty() bool {
	return sc.Len() == 0
}

// Clear removes all elements from the container.
func (sc *SliceContainer) Clear() {
	if sc.ThreadSafe {
		sc.Lock.Lock()
		defer sc.Lock.Unlock()
		sc.Backer = []ContainerElement{}
		return
	}

	sc.Backer = []ContainerElement{}
}

// Contains returns true if the given item is in the container.
func (sc *SliceContainer) Contains(item ContainerElement) bool {
	if sc.ThreadSafe {
		sc.Lock.Lock()
		defer sc.Lock.Unlock()
		return sc.containsHelper(item)
	}

	return sc.containsHelper(item)
}

// containsHelper returns whether the given element is in the container.
func (sc *SliceContainer) containsHelper(item ContainerElement) bool {
	for _, i := range sc.Backer {
		if i.Equals(item) {
			return true
		}
	}

	return false
}

// Add returns true if the given element was appended to the end of the container.
func (sc *SliceContainer) Add(item ContainerElement) bool {
	if sc.ThreadSafe {
		sc.Lock.Lock()
		defer sc.Lock.Unlock()
		sc.Backer = append(sc.Backer, item)
	} else {
		sc.Backer = append(sc.Backer, item)
	}

	return true
}

// Remove removes the element at the given index and
// returns whether the removal was successful.
func (sc *SliceContainer) Remove(item ContainerElement) bool {
	if sc.ThreadSafe {
		sc.Lock.Lock()
		defer sc.Lock.Unlock()
		idx := sc.findHelper(item)
		if idx < 0 {
			return false
		}
		sc.RemoveAtIndex(idx)
	} else {
		idx := sc.findHelper(item)
		if idx < 0 {
			return false
		}
		sc.RemoveAtIndex(idx)
	}

	return true
}

// findHelper searches the list for the given element and returns the index
// of the item in the container. If the item doesn't exist in the list, then -1
// is returned.
func (sc *SliceContainer) findHelper(item ContainerElement) int {
	for idx, val := range sc.Backer {
		if item.Equals(val) {
			return idx
		}
	}

	return -1
}

// RemoveAtIndex removes the element at the given idx.
func (sc *SliceContainer) RemoveAtIndex(idx int) {
	sc.Backer = append(sc.Backer[:idx], sc.Backer[idx+1:]...)

	// We should check to see if we need to resize the slice. We don't
	// want it to be the case that we added a ton of items then removed
	// a bunch and now we are still holding onto the large backing array
	// for the slice.
	emptyFactor := float32(len(sc.Backer)) / float32(cap(sc.Backer))
	if emptyFactor <= sc.ShrinkFactor {
		// Shrink the slice's capacity by 1/2
		newCap := cap(sc.Backer) / 2
		newBacker := make([]ContainerElement, len(sc.Backer), newCap)
		copy(newBacker, sc.Backer)
		sc.Backer = newBacker
	}
}
