package listadts

import (
	"sync"

	adts "github.com/johnsrd7/go-adts"
)

// SliceList is a simple type that implements the List interface.
type SliceList struct {
	backer       []adts.ContainerElement
	lock         *sync.Mutex
	threadSafe   bool
	shrinkFactor float32
}

// MakeSliceList creates a new non-threadsafe SliceList.
func MakeSliceList() *SliceList {
	return &SliceList{[]adts.ContainerElement{}, &sync.Mutex{}, false, 0.25}
}

// MakeSliceListThreadSafe creates a new threadsafe SliceList.
func MakeSliceListThreadSafe() *SliceList {
	return &SliceList{[]adts.ContainerElement{}, &sync.Mutex{}, true, 0.25}
}

// Cap returns the capacity of the list.
func (sl *SliceList) Cap() int {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return cap(sl.backer)
	}
	return cap(sl.backer)
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the list.
func (sl *SliceList) Len() int {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return len(sl.backer)
	}
	return len(sl.backer)
}

// IsEmpty returns if the list is empty or not.
func (sl *SliceList) IsEmpty() bool {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return sl.Len() == 0
	}

	return sl.Len() == 0
}

// Clear removes all elements from the list.
func (sl *SliceList) Clear() {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		sl.backer = []adts.ContainerElement{}
	}

	sl.backer = []adts.ContainerElement{}
}

// Contains returns true if the given item is in the list.
func (sl *SliceList) Contains(item adts.ContainerElement) bool {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return sl.containsHelper(item)
	}

	return sl.containsHelper(item)
}

// containsHelper returns whether the given element is in the list.
func (sl *SliceList) containsHelper(item adts.ContainerElement) bool {
	for _, i := range sl.backer {
		if i.Equals(item) {
			return true
		}
	}

	return false
}

// Add returns true if the given element was appended to the end of the list.
func (sl *SliceList) Add(item adts.ContainerElement) bool {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		sl.backer = append(sl.backer, item)
	} else {
		sl.backer = append(sl.backer, item)
	}

	return true
}

// Remove removes the element at the given index and
// returns whether the removal was successful.
func (sl *SliceList) Remove(item adts.ContainerElement) bool {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		idx := sl.findHelper(item)
		if idx < 0 {
			return false
		}
		sl.removeHelper(idx)
	} else {
		idx := sl.findHelper(item)
		if idx < 0 {
			return false
		}
		sl.removeHelper(idx)
	}

	return true
}

// findHelper searches the list for the given element and returns the index
// of the item in the list. If the item doesn't exist in the list, then -1
// is returned.
func (sl *SliceList) findHelper(item adts.ContainerElement) int {
	for idx, val := range sl.backer {
		if item.Equals(val) {
			return idx
		}
	}

	return -1
}

// removeHelper removes the element at the given idx.
func (sl *SliceList) removeHelper(idx int) {
	sl.backer = append(sl.backer[:idx], sl.backer[idx+1:]...)

	// We should check to see if we need to resize the slice. We don't
	// want it to be the case that we added a ton of items then removed
	// a bunch and now we are still holding onto the large backing array
	// for the slice.
	emptyFactor := float32(len(sl.backer)) / float32(cap(sl.backer))
	if emptyFactor <= sl.shrinkFactor {
		// Shrink the slice's capacity by 1/2
		newCap := cap(sl.backer) / 2
		newBacker := make([]adts.ContainerElement, len(sl.backer), newCap)
		copy(newBacker, sl.backer)
		sl.backer = newBacker
	}
}

// -------------------------------------------------------
// List Methods
// -------------------------------------------------------

// Get returns the element at the given index.
func (sl *SliceList) Get(idx int) adts.ContainerElement {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return sl.backer[idx]
	}

	return sl.backer[idx]
}

// Set changes the value at the given index to the given new value
// and returns the old value that was at the given index.
func (sl *SliceList) Set(idx int, newVal adts.ContainerElement) adts.ContainerElement {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		oldVal := sl.backer[idx]
		sl.backer[idx] = newVal
		return oldVal
	}

	oldVal := sl.backer[idx]
	sl.backer[idx] = newVal
	return oldVal
}
