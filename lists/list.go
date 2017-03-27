package listadts

import adts "github.com/johnsrd7/go-adts"

// List is common interface for a List ADT.
type List interface {
	Get(idx int) adts.ContainerElement
	Set(idx int, newVal adts.ContainerElement) adts.ContainerElement

	adts.Container
	// Len() int
	// IsEmpty() bool
	// Clear
	// Contains(item) bool
	// Add(item) bool
	// Remove(item) bool
}
