package listcontainer

import (
	"sync"
)

// SliceList is a simple type that implements the List interface.
type SliceList struct {
	backer       []ListElement
	lock         *sync.Mutex
	threadSafe   bool
	shrinkFactor float32
}

// MakeSliceList creates a new non-threadsafe SliceList.
func MakeSliceList() *SliceList {
	return &SliceList{[]ListElement{}, &sync.Mutex{}, false, 0.25}
}

// MakeSliceListThreadSafe creates a new threadsafe SliceList.
func MakeSliceListThreadSafe() *SliceList {
	return &SliceList{[]ListElement{}, &sync.Mutex{}, true, 0.25}
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

// Len returns the number of elements in the list.
func (sl *SliceList) Len() int {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return len(sl.backer)
	}
	return len(sl.backer)
}

// Contains returns whether the element is in the list.
func (sl *SliceList) Contains(item ListElement) bool {
	if sl.Len() > 0 &&
		sl.Get(0).TypeOf() != item.TypeOf() {
		return false
	}

	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return sl.containsHelper(item)
	}

	return sl.containsHelper(item)
}

// containsHelper returns whether the given element is in the list. Doesn't
// check to see if the element is of the proper type. It is assumed that
// the type check is done before the call to this method.
func (sl *SliceList) containsHelper(item ListElement) bool {
	for _, i := range sl.backer {
		if i.Equals(item) {
			return true
		}
	}

	return false
}

// Add returns true if the given element was appended to the end of the list.
func (sl *SliceList) Add(item ListElement) bool {
	if sl.Len() > 0 && sl.Get(0).TypeOf() != item.TypeOf() {
		return false
	}

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
func (sl *SliceList) Remove(idx int) bool {
	if idx >= sl.Len() {
		return false
	}

	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		sl.removeHelper(idx)
	} else {
		sl.removeHelper(idx)
	}

	return true
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
		newBacker := make([]ListElement, len(sl.backer), newCap)
		copy(newBacker, sl.backer)
		sl.backer = newBacker
	}
}

// Get returns the element at the given index.
func (sl *SliceList) Get(idx int) ListElement {
	if sl.threadSafe {
		sl.lock.Lock()
		defer sl.lock.Unlock()
		return sl.backer[idx]
	}

	return sl.backer[idx]
}
