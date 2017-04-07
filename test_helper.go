package adts

// IntElt is a wrapper for an int for testing
type IntElt int

// Equals returns true if the given container element is the same
// as this.
func (i IntElt) Equals(j ContainerElement) bool {
	if jElt, ok := j.(IntElt); ok {
		return i == jElt
	}
	return false
}
