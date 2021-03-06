package stackadts

import adts "github.com/johnsrd7/go-adts"

// Stack is a common interface for the Stack ADT.
type Stack interface {
	Push(adts.ContainerElement) bool
	Pop() (adts.ContainerElement, bool)

	adts.Container
	// Len() int
	// IsEmpty() bool
	// Clear()
	// Contains(item) bool
	// Add(item) bool
	// Remove(item) bool
}
