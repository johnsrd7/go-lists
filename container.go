package adts

// ContainerElement is a basic interface for any element that can be stored in a container.
type ContainerElement interface {
	Equals(ContainerElement) bool
}

// EmptyContainerElement is an empty container element that can be returned by
// any of the data structures for when an operation fails.
type EmptyContainerElement struct {
}

// Equals returns true if the given Container element is the same as this one.
func (ece EmptyContainerElement) Equals(other ContainerElement) bool {
	_, ok := other.(EmptyContainerElement)
	return ok
}

// Container is a basic interface for a container ADT.
type Container interface {
	Len() int
	IsEmpty() bool
	Clear()
	Contains(ContainerElement) bool
	Add(ContainerElement) bool
	Remove(ContainerElement) bool
}
