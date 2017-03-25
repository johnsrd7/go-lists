package listcontainer

import "sync"

// SliceList is a simple type that implements the List interface.
type SliceList struct {
	backer     []ListElement
	lock       *sync.Mutex
	threadSafe bool
}

// MakeSliceList creates a new non-threadsafe SliceList.
func MakeSliceList() *SliceList {
	return &SliceList{[]ListElement{}, &sync.Mutex{}, false}
}

// MakeSliceListThreadSafe creates a new threadsafe SliceList.
func MakeSliceListThreadSafe() *SliceList {
	return &SliceList{[]ListElement{}, &sync.Mutex{}, true}
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

	/*for i := idx; i < len(sl.backer)-1; i++ {
		sl.backer[i] = sl.backer[i+1]
	}

	sl.backer = sl.backer[:len(sl.backer)-1]*/
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
