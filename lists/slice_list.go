package listadts

import (
	adts "github.com/johnsrd7/go-adts"
)

// SliceList is a simple type that implements the List interface.
type SliceList struct {
	backer *adts.SliceContainer
}

// MakeSliceList creates a new non-threadsafe SliceList.
func MakeSliceList() *SliceList {
	return &SliceList{adts.MakeSliceContainer()}
}

// MakeSliceListThreadSafe creates a new threadsafe SliceList.
func MakeSliceListThreadSafe() *SliceList {
	return &SliceList{adts.MakeSliceContainerThreadSafe()}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the list.
func (sl *SliceList) Len() int {
	return sl.backer.Len()
}

// IsEmpty returns if the list is empty or not.
func (sl *SliceList) IsEmpty() bool {
	return sl.backer.Len() == 0
}

// Clear removes all elements from the list.
func (sl *SliceList) Clear() {
	sl.backer.Clear()
}

// Contains returns true if the given item is in the list.
func (sl *SliceList) Contains(item adts.ContainerElement) bool {
	return sl.backer.Contains(item)
}

// Add returns true if the given element was appended to the end of the list.
func (sl *SliceList) Add(item adts.ContainerElement) bool {
	return sl.backer.Add(item)
}

// Remove removes the element at the given index and
// returns whether the removal was successful.
func (sl *SliceList) Remove(item adts.ContainerElement) bool {
	return sl.backer.Remove(item)
}

// -------------------------------------------------------
// List Methods
// -------------------------------------------------------

// Get returns the element at the given index.
func (sl *SliceList) Get(idx int) adts.ContainerElement {
	if sl.backer.ThreadSafe {
		sl.backer.Lock.Lock()
		defer sl.backer.Lock.Unlock()
		return sl.backer.Backer[idx]
	}

	return sl.backer.Backer[idx]
}

// Set changes the value at the given index to the given new value
// and returns the old value that was at the given index.
func (sl *SliceList) Set(idx int, newVal adts.ContainerElement) adts.ContainerElement {
	if sl.backer.ThreadSafe {
		sl.backer.Lock.Lock()
		defer sl.backer.Lock.Unlock()
		oldVal := sl.backer.Backer[idx]
		sl.backer.Backer[idx] = newVal
		return oldVal
	}

	oldVal := sl.backer.Backer[idx]
	sl.backer.Backer[idx] = newVal
	return oldVal
}
