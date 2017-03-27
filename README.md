# GO ADTS
Go Implementations for data structures.

## Data Structures
- [Containers](#containers)
  - [Lists](#lists)
    - SliceList (Threadsafe and non-threadsafe)
	- SinglyLinkedList

## Containers
The following is the basic Container interface used by many of the data structures.
```go
type Container interface {
	Len() int
	IsEmpty() bool
	Clear()
	Contains(ContainerElement) bool
	Add(ContainerElement) bool
	Remove(ContainerElement) bool
}
```
The Container interface uses a basic ContainerElement interface for objects that it stores.
```go
type ContainerElement interface {
	Equals(ContainerElement) bool
}
```
