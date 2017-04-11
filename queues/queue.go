package queueadts

import adts "github.com/johnsrd7/go-adts"

// Queue is a common interface for the Queue ADT.
type Queue interface {
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
