package listcontainer

import "reflect"

// ListElement is a simple interface for comparing elements
type ListElement interface {
	TypeOf() reflect.Type
	Equals(ListElement) bool
}

// List is common interface for a List ADT.
type List interface {
	Len() int
	Contains(ListElement) bool
	Add(ListElement) bool
	Remove(int) bool
	Get(int) ListElement
}
