package listadts

import adts "github.com/johnsrd7/go-adts"

type intElt int

func (i intElt) Equals(j adts.ContainerElement) bool {
	if jElt, ok := j.(intElt); ok {
		return i == jElt
	}
	return false
}
