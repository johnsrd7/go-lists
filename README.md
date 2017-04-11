# GO ADTS
Go Implementations for data structures (both threadsafe and non-threadsafe).

**Currently threadsafe is untested. Use at your own risk.**

## Data Structures
- [Containers](#containers)
  - [Lists](#lists)
    - SliceList (Threadsafe and non-threadsafe)
	- SinglyLinkedList (Threadsafe and non-threadsafe)
  - [Stacks](#stacks)
    - SliceStack (Threadsafe and non-threadsafe)
	- ListStack (Threadsafe and non-threadsafe)

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

## Lists
The following is the basic List interface used by the list data structures.
```go
type List interface {
	Get(int) ContainerElement
	Set(int, ContainerElement) ContainerElement
	
	Container
}
```

## Stacks
The following is the basic Stack interface used by the stack data structures.
```go
type Stack interface {
	Push(ContainerElement)
	Pop() ContainerElement
	
	Container
}
```
