package adts

type ContainerElement interface {
	Equals(ContainerElement) bool
}

type Container interface {
	Len() int
	IsEmpty() bool
	Clear()
	Contains(ContainerElement) bool
	Add(ContainerElement) bool
	Remove(ContainerElement) bool
}
